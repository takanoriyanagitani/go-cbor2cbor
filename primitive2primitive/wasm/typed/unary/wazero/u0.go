package unary0

import (
	"context"
	"errors"
	"strconv"

	wa "github.com/tetratelabs/wazero/api"

	cnv "github.com/takanoriyanagitani/go-cbor2cbor/conv"

	pw "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive/wasm"
)

const (
	FunctionNameDefault string = "converter"
)

var (
	ErrInvalidResult   error = errors.New("invalid result")
	ErrInvalidFunction error = errors.New("invalid function")

	ErrInvalidInputType error = errors.New("invalid input type")

	ErrInvalidConverter error = errors.New("invalid converter")
)

type Func0 struct {
	wa.Function
}

var ConverterInvalid cnv.Converter = func(_ context.Context, _ any) (a any, e error) {
	return a, ErrInvalidConverter
}

var ConverterIdentity cnv.Converter = func(_ context.Context, input any) (any, error) {
	return input, nil
}

func (f Func0) ToConverter() cnv.Converter {
	var fd wa.FunctionDefinition = f.Function.Definition()

	var ptyps []wa.ValueType = fd.ParamTypes()
	if 1 != len(ptyps) {
		return ConverterInvalid
	}
	var ptyp wa.ValueType = ptyps[0]
	enc, ok := EncoderMap[ptyp]
	if !ok {
		return ConverterInvalid
	}

	var rtyps []wa.ValueType = fd.ResultTypes()
	if 1 != len(rtyps) {
		return ConverterInvalid
	}
	var rtyp wa.ValueType = rtyps[0]
	dec, ok := DecoderMap[rtyp]
	if !ok {
		return ConverterInvalid
	}

	return func(ctx context.Context, input any) (a any, e error) {
		encoded, e := enc(input)
		if nil != e {
			return a, e
		}

		results, e := f.Function.Call(ctx, encoded)
		if nil != e {
			return a, e
		}

		if 1 != len(results) {
			return a, ErrInvalidResult
		}

		var decoded any = dec(results[0])
		return decoded, nil
	}
}

type WasmIndexSource func(context.Context, uint32) ([]byte, error)

type IndexToName func(idx uint32) string

func IndexToNameDefault(idx uint32) string {
	return strconv.Itoa(int(idx))
}

type WasmExt string

var WasmExtDefault WasmExt = "wasm"

func (w WasmExt) IndexToBasenameNew(i2n IndexToName) IndexToName {
	var ext string = string(w)
	return func(idx uint32) string {
		var name string = i2n(idx)
		return name + "." + ext
	}
}

func (w WasmExt) ToWasmIndexSourceFS(
	i2n IndexToName,
	fssrc pw.WasmStore,
) WasmIndexSource {
	var idx2name IndexToName = w.IndexToBasenameNew(i2n)
	return func(ctx context.Context, idx uint32) ([]byte, error) {
		var basename string = idx2name(idx)
		return fssrc(ctx, basename)
	}
}
