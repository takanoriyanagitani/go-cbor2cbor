package arr2wz2arr

import (
	"context"

	w0 "github.com/tetratelabs/wazero"

	aa "github.com/takanoriyanagitani/go-cbor2cbor/arr2arr/arr2wasm2arr"
)

type Runtime struct {
	w0.Runtime

	w0.ModuleConfig

	SetInputSize       string
	EstimateOutputSize string
	SetOutputSize      string
	InputOffset        string
	Convert            string
	OutputOffset       string
}

func (r Runtime) Close(ctx context.Context) error {
	return r.Runtime.Close(ctx)
}

func (r Runtime) Compile(
	ctx context.Context,
	wasm []byte,
) (w0.CompiledModule, error) {
	return r.Runtime.CompileModule(ctx, wasm)
}

func (r Runtime) IntoCompiled(
	ctx context.Context,
	wasm []byte,
) (c Compiled, e error) {
	compiled, e := r.Compile(ctx, wasm)
	if nil != e {
		return c, e
	}

	c.Runtime = r.Runtime
	c.CompiledModule = compiled
	c.ModuleConfig = r.ModuleConfig

	c.SetInputSize = r.SetInputSize
	c.EstimateOutputSize = r.EstimateOutputSize
	c.SetOutputSize = r.SetOutputSize
	c.InputOffset = r.InputOffset
	c.Convert = r.Convert
	c.OutputOffset = r.OutputOffset

	return c, nil
}

func (r Runtime) WasmSourceIntoCompiled(
	ctx context.Context,
	ws aa.WasmSource,
) (c Compiled, e error) {
	wbytes, e := ws(ctx)
	if nil != e {
		return c, e
	}
	return r.IntoCompiled(ctx, wbytes)
}

func RuntimeNew(rtm w0.Runtime, cfg Config) Runtime {
	return Runtime{
		Runtime: rtm,

		ModuleConfig: cfg.ModuleConfig,

		SetInputSize:       cfg.SetInputSize,
		EstimateOutputSize: cfg.EstimateOutputSize,
		SetOutputSize:      cfg.SetOutputSize,
		InputOffset:        cfg.InputOffset,
		Convert:            cfg.Convert,
		OutputOffset:       cfg.OutputOffset,
	}
}

func RuntimeNewDefault(rtm w0.Runtime) Runtime {
	return RuntimeNew(rtm, ConfigDefault)
}
