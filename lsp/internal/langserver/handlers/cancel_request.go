package handlers

import (
	"context"
	"fmt"

	"github.com/creachadair/jrpc2"
	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
)

func CancelRequest(ctx context.Context, params lsp.CancelParams) error {
	id, err := decodeRequestID(params.ID)
	if err != nil {
		return err
	}

	jrpc2.ServerFromContext(ctx).CancelRequest(id)
	return nil
}

func decodeRequestID(v interface{}) (string, error) {
	if val, ok := v.(string); ok {
		return val, nil
	}
	if val, ok := v.(float64); ok {
		return fmt.Sprintf("%d", int64(val)), nil
	}

	return "", fmt.Errorf("unable to decode request ID: %#v", v)
}
