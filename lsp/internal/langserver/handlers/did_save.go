package handlers

import (
	"context"

	lsctx "github.com/glennsarti/blog-wasm-vscode/lsp/internal/context"
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
)

func (lh *logHandler) TextDocumentDidSave(ctx context.Context, params lsp.DidSaveTextDocumentParams) error {
	expFeatures, err := lsctx.ExperimentalFeatures(ctx)
	if err != nil {
		return err
	}
	if !expFeatures.ValidateOnSave {
		return nil
	}

	return err
}
