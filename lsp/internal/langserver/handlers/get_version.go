package handlers

import (
	"context"
	"errors"

	lsctx "github.com/glennsarti/blog-wasm-vscode/lsp/internal/context"
)

type VersionResponse struct {
	LspVersion string `json:"lspVersion"`
}

func DemoGetVersion(ctx context.Context) (VersionResponse, error) {
	if ver, ok := lsctx.LanguageServerVersion(ctx); ok {
		// Quick and dirty response
		return VersionResponse{LspVersion: ver}, nil
	} else {
		return VersionResponse{LspVersion: ""}, errors.New("Could not determine Language Server Version")
	}
}
