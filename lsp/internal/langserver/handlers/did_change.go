package handlers

import (
	"context"
	"fmt"

	lsctx "github.com/glennsarti/blog-wasm-vscode/lsp/internal/context"
	demo "github.com/glennsarti/blog-wasm-vscode/lsp/internal/langserver/demo"
	ilsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/lsp"
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
)

func TextDocumentDidChange(ctx context.Context, params lsp.DidChangeTextDocumentParams) error {
	p := lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{
				URI: params.TextDocument.URI,
			},
			Version: params.TextDocument.Version,
		},
		ContentChanges: params.ContentChanges,
	}

	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return err
	}

	fh := ilsp.VersionedFileHandler(p.TextDocument)
	f, err := fs.GetDocument(fh)
	if err != nil {
		return err
	}

	// Versions don't have to be consecutive, but they must be increasing
	if int(p.TextDocument.Version) <= f.Version() {
		fs.CloseAndRemoveDocument(fh)
		return fmt.Errorf("Old version (%d) received, current version is %d. "+
			"Unable to update %s. This is likely a bug, please report it.",
			int(p.TextDocument.Version), f.Version(), p.TextDocument.URI)
	}

	changes, err := ilsp.DocumentChanges(params.ContentChanges, f)
	if err != nil {
		return err
	}
	err = fs.ChangeDocument(fh, changes)
	if err != nil {
		return err
	}

	if doc, err := fs.GetDocument(fh); err != nil {
		return err
	} else {
		demo.GenerateAndPublishDiagnostics(ctx, doc)
	}

	return nil
}
