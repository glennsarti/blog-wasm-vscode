package handlers

import (
	"context"

	lsctx "github.com/glennsarti/blog-wasm-vscode/lsp/internal/context"
	ilsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/lsp"
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
)

func TextDocumentDidClose(ctx context.Context, params lsp.DidCloseTextDocumentParams) error {
	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return err
	}

	fh := ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI)
	err = fs.CloseAndRemoveDocument(fh)
	if err != nil {
		return err
	}

	return nil
}
