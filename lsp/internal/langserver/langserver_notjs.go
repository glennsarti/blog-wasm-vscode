//go:build !js

package langserver

import "errors"

func (ls *langServer) StartWASM() error {
	return errors.New("not compiled for WASM")
}
