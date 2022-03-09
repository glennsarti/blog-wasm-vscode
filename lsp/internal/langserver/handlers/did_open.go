package handlers

import (
	"context"

	lsctx "github.com/glennsarti/blog-wasm-vscode/lsp/internal/context"
	ilsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/lsp"
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"

	demo "github.com/glennsarti/blog-wasm-vscode/lsp/internal/langserver/demo"
)

func (lh *logHandler) TextDocumentDidOpen(ctx context.Context, params lsp.DidOpenTextDocumentParams) error {
	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return err
	}

	f := ilsp.FileFromDocumentItem(params.TextDocument)
	err = fs.CreateAndOpenDocument(f, f.LanguageID(), f.Text())
	if err != nil {
		return err
	}

	if doc, err := fs.GetDocument(f); err != nil {
		return err
	} else {
		demo.GenerateAndPublishDiagnostics(ctx, doc)
	}

	return nil
}
