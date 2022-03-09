package demo

import (
	"context"
	"strings"

	lsctx "github.com/glennsarti/blog-wasm-vscode/lsp/internal/context"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/filesystem"
	"github.com/glennsarti/blog-wasm-vscode/lsp/internal/langserver/diagnostics"
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
)

const LspDiagnosticSource = "demo"

func GenerateAndPublishDiagnostics(ctx context.Context, f filesystem.Document) (err error) {
	list, err := generate(f)
	if err != nil {
		return err
	}
	notifier, err := lsctx.DiagnosticsNotifier(ctx)
	if err != nil {
		return err
	}

	notifier.PublishDiagnostics(ctx, "", list)
	return nil
}

func generate(f filesystem.Document) (list diagnostics.Diagnostics, err error) {
	diags := make(diagnostics.Diagnostics, 0)
	src := diagnostics.DiagnosticSource(f.FullPath())
	diags[src] = make(diagnostics.DiagnosticList, 0)

	for lineNum, line := range f.Lines() {
		text := string(line.Bytes())
		if strings.Contains(text, "demo") {
			diag := lsp.Diagnostic{
				Message: "Line contains the word demo",
				Source:  LspDiagnosticSource,
				Range: lsp.Range{
					Start: lsp.Position{
						Line:      uint32(lineNum),
						Character: 0,
					},
					End: lsp.Position{
						Line:      uint32(lineNum),
						Character: uint32(len(text)),
					},
				},
				Severity: lsp.SeverityError,
			}
			diags[src] = append(diags[src], diag)
		}
	}

	return diags, nil
}
