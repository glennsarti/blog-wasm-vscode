package handlers

import (
	"context"
	"path/filepath"
	"strings"

	jrpc2 "github.com/creachadair/jrpc2"
	lsctx "github.com/glennsarti/blog-wasm-vscode/lsp/internal/context"
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
	"github.com/mitchellh/go-homedir"
)

func (svc *service) Initialize(ctx context.Context, params lsp.InitializeParams) (lsp.InitializeResult, error) {
	serverCaps := lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: lsp.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    lsp.Full,
			},
			DefinitionProvider:         false,
			ReferencesProvider:         false,
			HoverProvider:              false,
			DocumentFormattingProvider: false,
			DocumentSymbolProvider:     false,
			WorkspaceSymbolProvider:    false,
			Workspace: lsp.Workspace5Gn{
				WorkspaceFolders: lsp.WorkspaceFolders4Gn{
					Supported:           false,
					ChangeNotifications: "workspace/didChangeWorkspaceFolders",
				},
			},
		},
	}

	var err error

	serverCaps.ServerInfo.Name = "demo-ls"
	version, ok := lsctx.LanguageServerVersion(ctx)
	if ok {
		serverCaps.ServerInfo.Version = version
	}

	svc.server = jrpc2.ServerFromContext(ctx)

	err = svc.configureSessionDependencies(ctx)
	if err != nil {
		return serverCaps, err
	}

	return serverCaps, nil
}

func resolvePath(rootDir, rawPath string) (string, error) {
	path, err := homedir.Expand(rawPath)
	if err != nil {
		return "", err
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(rootDir, rawPath)
	}

	return cleanupPath(path)
}

func cleanupPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	return toLowerVolumePath(absPath), err
}

func toLowerVolumePath(path string) string {
	volume := filepath.VolumeName(path)
	return strings.ToLower(volume) + path[len(volume):]
}
