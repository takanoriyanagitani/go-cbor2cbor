package unary0

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"maps"

	w0 "github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"

	cnv "github.com/takanoriyanagitani/go-cbor2cbor/conv"

	pw "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive/wasm"
)

type Modules struct {
	Indexed      map[uint32]wa.Module
	FunctionName string
}

func (m Modules) CloseAll(ctx context.Context) error {
	var err []error
	var modules iter.Seq[wa.Module] = maps.Values(m.Indexed)
	for mdl := range modules {
		err = append(err, mdl.Close(ctx))
	}
	return errors.Join(err...)
}

func (m Modules) ToConverterMap(ctx context.Context) (map[uint32]cnv.Converter, error) {
	ret := map[uint32]cnv.Converter{}
	var i iter.Seq2[uint32, wa.Module] = maps.All(m.Indexed)
	for key, val := range i {
		var f wa.Function = val.ExportedFunction(m.FunctionName)
		if nil == f {
			return nil, fmt.Errorf(
				"%w: function not found(name=%s)",
				ErrInvalidFunction,
				m.FunctionName,
			)
		}
		ret[key] = Func0{f}.ToConverter()
	}
	return ret, nil
}

type ModuleSource func(context.Context, []byte) (wa.Module, error)

func ModuleSourceNew(rtm w0.Runtime, cfg w0.ModuleConfig) ModuleSource {
	return func(ctx context.Context, wasm []byte) (wa.Module, error) {
		return rtm.InstantiateWithConfig(ctx, wasm, cfg)
	}
}

func ModuleSourceNewDefault(rtm w0.Runtime) ModuleSource {
	return ModuleSourceNew(rtm, ModuleConfigDefault)
}

func ModulesNew(
	ctx context.Context,
	indices []uint32,
	funcName string,
	wasmSource WasmIndexSource,
	instantiate ModuleSource,
) (Modules, error) {
	m := Modules{
		Indexed:      map[uint32]wa.Module{},
		FunctionName: funcName,
	}

	for _, idx := range indices {
		source, e := wasmSource(ctx, idx)
		if nil != e {
			return Modules{}, errors.Join(e, m.CloseAll(ctx))
		}

		mdl, e := instantiate(ctx, source)
		if nil != e {
			return Modules{}, errors.Join(e, m.CloseAll(ctx))
		}
		m.Indexed[idx] = mdl
	}
	return m, nil
}

func ModulesNewDefaultFS(
	ctx context.Context,
	indices []uint32,
	fssrc pw.WasmStore,
	rtm w0.Runtime,
) (Modules, error) {
	return ModulesNew(
		ctx,
		indices,
		FunctionNameDefault,
		WasmExtDefault.ToWasmIndexSourceFS(
			IndexToNameDefault,
			fssrc,
		),
		ModuleSourceNewDefault(rtm),
	)
}

var ModuleConfigDefault w0.ModuleConfig = w0.NewModuleConfig().
	WithName("")
