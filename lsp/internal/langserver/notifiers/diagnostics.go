package notifiers

import (
	"context"
	"log"
	"path/filepath"
	"sync"

	"github.com/creachadair/jrpc2"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/langserver/diagnostics"
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/uri"
)

type diagContext struct {
	ctx   context.Context
	uri   lsp.DocumentURI
	diags []lsp.Diagnostic
}

// Notifier is a type responsible for queueing diagnostics to be converted
// and sent to the client
type DiagnosticsNotifier struct {
	logger *log.Logger
	// sessCtx        context.Context
	diags          chan diagContext
	closeDiagsOnce sync.Once
}

func NewDiagnosticsNotifier(logger *log.Logger) *DiagnosticsNotifier {
	n := &DiagnosticsNotifier{
		logger: logger,
		//sessCtx: sessCtx,
		diags: make(chan diagContext, 50),
	}
	go n.notify()
	return n
}

// PublishHCLDiags accepts a map of HCL diagnostics per file and queues them for publishing.
// A dir path is passed which is joined with the filename keys of the map, to form a file URI.
func (n *DiagnosticsNotifier) PublishDiagnostics(ctx context.Context, dirPath string, diags diagnostics.Diagnostics) {
	select {
	case <-ctx.Done():
		n.closeDiagsOnce.Do(func() {
			close(n.diags)
		})
		return
	default:
	}

	for filename, ds := range diags {
		n.diags <- diagContext{
			ctx:   ctx,
			uri:   lsp.DocumentURI(uri.FromPath(filepath.Join(dirPath, string(filename)))),
			diags: ds,
		}
	}
}

func (n *DiagnosticsNotifier) notify() {
	for d := range n.diags {
		if err := jrpc2.ServerFromContext(d.ctx).Notify(d.ctx, "textDocument/publishDiagnostics", lsp.PublishDiagnosticsParams{
			URI:         d.uri,
			Diagnostics: d.diags,
		}); err != nil {
			n.logger.Printf("Error pushing diagnostics: %s", err)
		}
	}
}
