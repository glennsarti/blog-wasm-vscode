package lsp

import (
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
	"github.com/hashicorp/hcl-lang/lang"
)

func markupContent(content lang.MarkupContent, mdSupported bool) lsp.MarkupContent {
	return lsp.MarkupContent{
		Kind:  lsp.PlainText,
		Value: "ERROR!!",
	}

	// value := content.Value

	// kind := lsp.PlainText
	// if content.Kind == lang.MarkdownKind {
	// 	if mdSupported {
	// 		kind = lsp.Markdown
	// 	} else {
	// 		value = mdplain.Clean(value)
	// 	}
	// }

	// return lsp.MarkupContent{
	// 	Kind:  kind,
	// 	Value: value,
	// }
}
