package diagnostics

import (
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
)

// type DiagnosticsGenerator interface {
// 	Generate() (diags Diagnostics, err error)
// }

// Diagnostics is a list of Diagnostic instances.
type DiagnosticList []lsp.Diagnostic

type Diagnostics map[DiagnosticSource]DiagnosticList

type DiagnosticSource string

func NewDiagnostics() Diagnostics {
	return make(Diagnostics, 0)
}
