package handlers

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/code"
	rpch "github.com/creachadair/jrpc2/handler"

	lsctx "github.com/glennsarti/blog-wasm-vscode/lsp/internal/context"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/filesystem"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/langserver/notifiers"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/langserver/session"
	ilsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/lsp"
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/settings"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/state"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/telemetry"
	"github.com/hashicorp/hcl-lang/decoder"
	"github.com/hashicorp/hcl-lang/lang"
)

type service struct {
	logger *log.Logger

	srvCtx context.Context

	sessCtx     context.Context
	stopSession context.CancelFunc

	fs            filesystem.Filesystem
	telemetry     telemetry.Sender
	hclDecoder    *decoder.Decoder
	stateStore    *state.StateStore
	server        session.Server
	diagsNotifier *notifiers.DiagnosticsNotifier

	additionalHandlers map[string]rpch.Func
}

var discardLogs = log.New(ioutil.Discard, "", 0)

func NewSession(srvCtx context.Context) session.Session {
	fs := filesystem.NewFilesystem()

	sessCtx, stopSession := context.WithCancel(srvCtx)
	return &service{
		logger:      discardLogs,
		fs:          fs,
		srvCtx:      srvCtx,
		sessCtx:     sessCtx,
		stopSession: stopSession,
		telemetry:   &telemetry.NoopSender{},
	}
}

func (svc *service) SetLogger(logger *log.Logger) {
	svc.logger = logger
}

// Assigner builds out the jrpc2.Map according to the LSP protocol
// and passes related dependencies to handlers via context
func (svc *service) Assigner() (jrpc2.Assigner, error) {
	svc.logger.Println("Preparing new session ...")

	session := session.NewSession(svc.stopSession)

	err := session.Prepare()
	if err != nil {
		return nil, fmt.Errorf("Unable to prepare session: %w", err)
	}

	svc.telemetry = &telemetry.NoopSender{Logger: svc.logger}
	svc.fs.SetLogger(svc.logger)

	lh := LogHandler(svc.logger)
	cc := &lsp.ClientCapabilities{}

	rootDir := ""
	clientName := ""
	var expFeatures settings.ExperimentalFeatures

	m := map[string]rpch.Func{
		"initialize": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			err := session.Initialize(req)
			if err != nil {
				return nil, err
			}

			ctx = ilsp.WithClientCapabilities(ctx, cc)
			ctx = lsctx.WithRootDirectory(ctx, &rootDir)
			ctx = ilsp.ContextWithClientName(ctx, &clientName)
			ctx = lsctx.WithExperimentalFeatures(ctx, &expFeatures)

			version, ok := lsctx.LanguageServerVersion(svc.srvCtx)
			if ok {
				ctx = lsctx.WithLanguageServerVersion(ctx, version)
			}

			return handle(ctx, req, svc.Initialize)
		},
		"initialized": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			err := session.ConfirmInitialization(req)
			if err != nil {
				return nil, err
			}

			return handle(ctx, req, Initialized)
		},
		"textDocument/didChange": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			err := session.CheckInitializationIsConfirmed()
			if err != nil {
				return nil, err
			}
			ctx = lsctx.WithDocumentStorage(ctx, svc.fs)
			ctx = lsctx.WithDiagnosticsNotifier(ctx, svc.diagsNotifier)
			return handle(ctx, req, TextDocumentDidChange)
		},
		"textDocument/didOpen": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			err := session.CheckInitializationIsConfirmed()
			if err != nil {
				return nil, err
			}
			ctx = lsctx.WithDocumentStorage(ctx, svc.fs)
			ctx = lsctx.WithDiagnosticsNotifier(ctx, svc.diagsNotifier)
			return handle(ctx, req, lh.TextDocumentDidOpen)
		},
		"textDocument/didClose": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			err := session.CheckInitializationIsConfirmed()
			if err != nil {
				return nil, err
			}
			ctx = lsctx.WithDocumentStorage(ctx, svc.fs)
			return handle(ctx, req, TextDocumentDidClose)
		},

		"textDocument/didSave": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			err := session.CheckInitializationIsConfirmed()
			if err != nil {
				return nil, err
			}

			ctx = lsctx.WithDiagnosticsNotifier(ctx, svc.diagsNotifier)
			ctx = lsctx.WithExperimentalFeatures(ctx, &expFeatures)

			return handle(ctx, req, lh.TextDocumentDidSave)
		},

		// Can't disable these yet :-(
		"textDocument/codeLens": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			return nil, nil
		},
		"textDocument/documentLink": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			return nil, nil
		},

		// Custom messages
		"demo/getVersion": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			return DemoGetVersion(svc.srvCtx)
		},

		"shutdown": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			err := session.Shutdown(req)
			if err != nil {
				return nil, err
			}
			ctx = lsctx.WithDocumentStorage(ctx, svc.fs)
			svc.shutdown()
			return handle(ctx, req, Shutdown)
		},
		"exit": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			err := session.Exit()
			if err != nil {
				return nil, err
			}

			svc.stopSession()

			return nil, nil
		},
		"$/cancelRequest": func(ctx context.Context, req *jrpc2.Request) (interface{}, error) {
			err := session.CheckInitializationIsConfirmed()
			if err != nil {
				return nil, err
			}

			return handle(ctx, req, CancelRequest)
		},
	}

	// For use in tests, e.g. to test request cancellation
	if len(svc.additionalHandlers) > 0 {
		for methodName, handlerFunc := range svc.additionalHandlers {
			m[methodName] = handlerFunc
		}
	}

	return convertMap(m), nil
}

func (svc *service) configureSessionDependencies(_ context.Context) error {
	svc.diagsNotifier = notifiers.NewDiagnosticsNotifier(svc.logger)

	svc.stateStore.SetLogger(svc.logger)

	svc.hclDecoder = decoder.NewDecoder(nil) // TODO: Is this right?

	return nil
}

func (svc *service) setupTelemetry(version int, notifier session.ClientNotifier) error {
	t, err := telemetry.NewSender(version, notifier)
	if err != nil {
		return err
	}

	svc.telemetry = t
	return nil
}

func (svc *service) Finish(_ jrpc2.Assigner, status jrpc2.ServerStatus) {
	if status.Closed || status.Err != nil {
		svc.logger.Printf("session stopped unexpectedly (err: %v)", status.Err)
	}

	svc.shutdown()
	svc.stopSession()
}

func (svc *service) shutdown() {
}

// convertMap is a helper function allowing us to omit the jrpc2.Func
// signature from the method definitions
func convertMap(m map[string]rpch.Func) rpch.Map {
	hm := make(rpch.Map, len(m))

	for method, fun := range m {
		hm[method] = rpch.New(fun)
	}

	return hm
}

const requestCancelled code.Code = -32800

// handle calls a jrpc2.Func compatible function
func handle(ctx context.Context, req *jrpc2.Request, fn interface{}) (interface{}, error) {
	f := rpch.New(fn)
	result, err := f.Handle(ctx, req)
	if ctx.Err() != nil && errors.Is(ctx.Err(), context.Canceled) {
		err = fmt.Errorf("%w: %s", requestCancelled.Err(), err)
	}
	return result, err
}

func (svc *service) hclDecoderForDocument(ctx context.Context, doc filesystem.Document) (*decoder.PathDecoder, error) {
	return svc.hclDecoder.Path(lang.Path{
		Path:       doc.Dir(),
		LanguageID: doc.LanguageID(),
	})
}
