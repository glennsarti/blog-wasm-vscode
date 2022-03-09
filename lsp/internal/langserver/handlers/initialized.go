package handlers

import (
	"context"

	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
)

func Initialized(ctx context.Context, params lsp.InitializedParams) error {
	return nil
}
