package arr2wz2arr

import (
	"context"
	"errors"

	w0 "github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"
)

type Compiled struct {
	w0.Runtime
	w0.CompiledModule
	w0.ModuleConfig

	SetInputSize       string
	EstimateOutputSize string
	SetOutputSize      string
	InputOffset        string
	Convert            string
	OutputOffset       string
}

func (c Compiled) Close(ctx context.Context) error {
	return errors.Join(c.CompiledModule.Close(ctx), c.Runtime.Close(ctx))
}

func (c Compiled) InstantiateModule(ctx context.Context) (wa.Module, error) {
	return c.Runtime.InstantiateModule(
		ctx,
		c.CompiledModule,
		c.ModuleConfig,
	)
}

func (c Compiled) Instantiate(ctx context.Context) (conv Converter, e error) {
	mdl, e := c.InstantiateModule(ctx)
	if nil != e {
		return conv, e
	}
	conv.Module = mdl

	conv.Memory = mdl.Memory()

	conv.SetInputSize.Function = mdl.ExportedFunction(c.SetInputSize)
	conv.EstimateOutputSize.Function = mdl.ExportedFunction(c.EstimateOutputSize)
	conv.SetOutputSize.Function = mdl.ExportedFunction(c.SetOutputSize)
	conv.InputOffset.Function = mdl.ExportedFunction(c.InputOffset)
	conv.Convert.Function = mdl.ExportedFunction(c.Convert)
	conv.OutputOffset.Function = mdl.ExportedFunction(c.OutputOffset)

	e = conv.Validate()
	if nil != e {
		return Converter{}, errors.Join(e, conv.Close(ctx))
	}

	return conv, nil
}
