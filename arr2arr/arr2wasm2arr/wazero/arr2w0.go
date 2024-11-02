package arr2wz2arr

import (
	"context"
	"errors"
	"fmt"

	wa "github.com/tetratelabs/wazero/api"

	aa "github.com/takanoriyanagitani/go-cbor2cbor/arr2arr"
)

var (
	ErrUnableToWriteToMemory error = errors.New("unable to write to memory")
	ErrUnableToViewMemory    error = errors.New("unable to view memory")

	ErrInvalidModule error = errors.New("invalid module")

	ErrUnexpectedInteger error = errors.New("negative integer got")
)

type Convert struct {
	wa.Function
}

func IntToUintNonNeg(i int32, e error) (uint32, error) {
	if nil != e {
		return 0, e
	}
	if i < 0 {
		return 0, ErrUnexpectedInteger
	}
	return uint32(i), nil
}

func (c Convert) ConvertCbor(ctx context.Context) (uint32, error) {
	results, e := c.Function.Call(ctx)
	result := Result{
		Err: e,
		Raw: results,
	}
	i, e := result.GetInt32i(ctx)
	return IntToUintNonNeg(i, e)
}

type Converter struct {
	wa.Module
	wa.Memory

	SetInputSize
	EstimateOutputSize
	SetOutputSize
	InputOffset
	Convert
	OutputOffset
}

func (c Converter) Validate() error {
	oks := []bool{
		nil != c.Module,
		nil != c.Memory,

		nil != c.SetInputSize.Function,
		nil != c.EstimateOutputSize.Function,
		nil != c.SetOutputSize.Function,
		nil != c.InputOffset.Function,
		nil != c.Convert.Function,
		nil != c.OutputOffset.Function,
	}

	for _, ok := range oks {
		var ng bool = !ok
		if ng {
			return ErrInvalidModule
		}
	}
	return nil
}

func (c Converter) Close(ctx context.Context) error {
	return c.Module.Close(ctx)
}

func (c Converter) CopyToWasm(offset uint32, data []byte) error {
	var ok bool = c.Memory.Write(offset, data)
	switch ok {
	case true:
		return nil
	default:
		return ErrUnableToWriteToMemory
	}
}

func (c Converter) WasmToSlice(
	offset uint32,
	count uint32,
) (view []byte, e error) {
	view, ok := c.Memory.Read(offset, count)
	switch ok {
	case true:
		return view, nil
	default:
		return nil, ErrUnableToViewMemory
	}
}

func (c Converter) ConvertCbor(
	ctx context.Context,
	inputCbor []byte,
) (outputCbor []byte, e error) {
	var isz uint32 = uint32(len(inputCbor))
	_, e = c.SetInputSize.SetInputDefault(ctx, isz)
	if nil != e {
		return nil, fmt.Errorf("%w: unable to set input size", e)
	}

	osz, e := c.EstimateOutputSize.Estimate(ctx)
	if nil != e {
		return nil, fmt.Errorf("%w: unable to estimate output size", e)
	}

	_, e = c.SetOutputSize.SetOutSize(ctx, osz)
	if nil != e {
		return nil, fmt.Errorf("%w: unable to set output size", e)
	}

	ioff, e := c.InputOffset.GetInputOffset(ctx)
	if nil != e {
		return nil, fmt.Errorf("%w: unable to get input offset", e)
	}

	e = c.CopyToWasm(ioff, inputCbor)
	if nil != e {
		return nil, fmt.Errorf("%w: unable to copy to wasm", e)
	}

	csz, e := c.Convert.ConvertCbor(ctx)
	if nil != e {
		return nil, fmt.Errorf("%w: unable to convert", e)
	}

	ooff, e := c.OutputOffset.GetOutputOffset(ctx)
	if nil != e {
		return nil, fmt.Errorf("%w: unable to get output offset", e)
	}

	view, e := c.WasmToSlice(ooff, csz)
	if nil != e {
		return nil, fmt.Errorf("%w: unable to get view", e)
	}

	return view, nil
}

func (c Converter) ToRawArrToRawArr() aa.RawArrayToRawArray {
	return func(ctx context.Context, i aa.InputCbor) (aa.OutputCbor, error) {
		return c.ConvertCbor(ctx, i)
	}
}
