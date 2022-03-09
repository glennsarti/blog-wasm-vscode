package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/langserver"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/langserver/handlers"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/logging"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/pathtpl"

	lsctx "github.com/glennsarti/blog-wasm-vscode/lsp/internal/context"
)

type ServeCommandRunner struct {
	Version string
	// flags
	port           int
	logFilePath    string
	cpuProfile     string
	memProfile     string
	reqConcurrency int
}

func (c *ServeCommandRunner) flags() *flag.FlagSet {
	fs := defaultFlagSet("serve")

	fs.IntVar(&c.port, "port", 30337, "port number to listen on (turns server into TCP mode)")
	fs.StringVar(&c.logFilePath, "log-file", "", "path to a file to log into with support "+
		"for variables (e.g. Timestamp, Pid, Ppid) via Go template syntax {{.VarName}}")

	return fs
}

func (c *ServeCommandRunner) Run(args []string) (int, error) {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		return 1, errors.New(fmt.Sprintf("Error parsing command-line flags: %s", err))
	}

	if c.cpuProfile != "" {
		stop, err := writeCpuProfileInto(c.cpuProfile)
		defer stop()
		if err != nil {
			return 1, errors.New(err.Error())
		}
	}

	if c.memProfile != "" {
		defer writeMemoryProfileInto(c.memProfile)
	}

	var logger *log.Logger
	if c.logFilePath != "" {
		fl, err := logging.NewFileLogger(c.logFilePath)
		if err != nil {
			return 1, errors.New(fmt.Sprintf("Failed to setup file logging: %s", err))
		}
		defer fl.Close()

		logger = fl.Logger()
	} else {
		logger = logging.NewLogger(os.Stderr)
	}

	ctx, cancelFunc := lsctx.WithSignalCancel(context.Background(), logger,
		os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()
	//TODO:  DISABLE FOR GOOS
	c.reqConcurrency = 1
	if c.reqConcurrency != 0 {
		ctx = langserver.WithRequestConcurrency(ctx, c.reqConcurrency)
		logger.Printf("Custom request concurrency set to %d", c.reqConcurrency)
	}

	logger.Printf("Starting demo-ls %s", c.Version)

	ctx = lsctx.WithLanguageServerVersion(ctx, c.Version)

	srv := langserver.NewLangServer(ctx, handlers.NewSession)
	srv.SetLogger(logger)

	if runtime.GOOS == "js" {
		err := srv.StartWASM()
		if err != nil {
			return 1, errors.New(fmt.Sprintf("Failed to start StartWASM server: %s", err))
		}
		return 0, nil
	}

	if c.port != 0 {
		err := srv.StartTCP(fmt.Sprintf("localhost:%d", c.port))
		if err != nil {
			return 1, errors.New(fmt.Sprintf("Failed to start TCP server: %s", err))
		}
		return 0, nil
	}

	err := srv.StartAndWait(os.Stdin, os.Stdout)
	if err != nil {
		return 1, errors.New(fmt.Sprintf("Failed to start server: %s", err))
	}

	return 0, nil
}

type stopFunc func() error

func writeCpuProfileInto(rawPath string) (stopFunc, error) {
	path, err := pathtpl.ParseRawPath("cpuprofile-path", rawPath)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("could not create CPU profile: %s", err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		return f.Close, fmt.Errorf("could not start CPU profile: %s", err)
	}

	return func() error {
		pprof.StopCPUProfile()
		return f.Close()
	}, nil
}

func writeMemoryProfileInto(rawPath string) error {
	path, err := pathtpl.ParseRawPath("memprofile-path", rawPath)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create memory profile: %s", err)
	}
	defer f.Close()

	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		return fmt.Errorf("could not write memory profile: %s", err)
	}

	return nil
}
