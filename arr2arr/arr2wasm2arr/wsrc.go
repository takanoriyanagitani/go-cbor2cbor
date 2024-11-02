package arr2wasm

import (
	"context"
)

type WasmSource func(context.Context) (wasmBytes []byte, e error)
