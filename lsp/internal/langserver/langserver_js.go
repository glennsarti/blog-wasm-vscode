package langserver

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall/js"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

type wasmServer struct {
	logger        *log.Logger
	client        *jrpc2.Client
	ClientChannel channel.Channel
	ServerChannel channel.Channel
}

func (w *wasmServer) lsponmessage(this js.Value, inputs []js.Value) interface{} {
	message := inputs[0].String()
	w.ClientChannel.Send([]byte(message))
	return nil
}

func (w *wasmServer) Start() error {
	// Register WASM functions
	js.Global().Set("lsponmessage", js.FuncOf(w.lsponmessage))

	c, s := channel.Direct()
	w.ClientChannel = c
	w.ServerChannel = s

	w.asyncReader()

	return nil
}

func (w *wasmServer) asyncReader() {
	go func() error {
		w.logger.Printf("WASM Client Reader Started")
		for {
			buf, _ := w.ClientChannel.Recv()
			// TODO: Need a select here to bail when an error occurs
			js.Global().Call("goLangPostMessage", string(buf))
		}
		return nil
	}()
}

func (w *wasmServer) AssertReady() {
	w.logger.Printf("Flagging the WASM server is ready...")
	js.Global().Call("workerReady")
}

func newWasmServer(srvCtx context.Context, logger *log.Logger) *wasmServer {
	return &wasmServer{
		logger: logger,
	}
}

// Need a singleton here
var wasm *wasmServer

func (ls *langServer) StartWASM() error {
	wasm = newWasmServer(ls.srvCtx, ls.logger)

	err := wasm.Start()
	if err != nil {
		return fmt.Errorf("WASM services failed to start: %s", err)
	}
	ls.logger.Printf("WASM services started")

	srv, err := Server(ls.newService(), ls.srvOptions)
	if err != nil {
		return err
	}

	ls.logger.Printf("WASM LSP server starting up...")
	srv.Start(wasm.ServerChannel)
	ls.logger.Printf("WASM LSP server started")

	// Wrap waiter with a context so that we can cancel it here
	// after the service is cancelled (and srv.Wait returns)
	ctx, cancelFunc := context.WithCancel(ls.srvCtx)
	go func() {
		srv.Wait()
		cancelFunc()
	}()

	wasm.AssertReady()
	select {
	case <-ctx.Done():
		ls.logger.Printf("Stopping server (pid %d) ...", os.Getpid())
		srv.Stop()
	}

	ls.logger.Printf("Server (pid %d) stopped.", os.Getpid())
	return nil
}
