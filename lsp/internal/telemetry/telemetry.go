package telemetry

import (
	"context"
	"fmt"

	lsp "github.com/glennsarti/blog-wasm-vscode/lsp/internal/protocol"
)

type Telemetry struct {
	version  int
	notifier Notifier
}

type Notifier interface {
	Notify(ctx context.Context, method string, params interface{}) error
}

func NewSender(version int, notifier Notifier) (*Telemetry, error) {
	if version != lsp.TelemetryFormatVersion {
		return nil, fmt.Errorf("unsupported telemetry format version: %d", version)
	}

	return &Telemetry{
		version:  version,
		notifier: notifier,
	}, nil
}

func (t *Telemetry) SendEvent(ctx context.Context, name string, properties map[string]interface{}) {
	t.notifier.Notify(ctx, "telemetry/event", lsp.TelemetryEvent{
		Version:    t.version,
		Name:       name,
		Properties: properties,
	})
}
