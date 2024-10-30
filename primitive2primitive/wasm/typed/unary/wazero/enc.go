package unary0

import (
	wa "github.com/tetratelabs/wazero/api"

	util "github.com/takanoriyanagitani/go-cbor2cbor/util"

	wt "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive/wasm/typed"
)

type InputEncoder func(any) (uint64, error)

var EncodeI64 InputEncoder = util.ComposeErr(
	wt.AnyToI64,
	util.OkFunc(wa.EncodeI64),
)

var EncodeI32 InputEncoder = util.ComposeErr(
	wt.AnyToI32,
	util.OkFunc(wa.EncodeI32),
)

var EncodeF32 InputEncoder = util.ComposeErr(
	wt.AnyToF32,
	util.OkFunc(wa.EncodeF32),
)

var EncodeF64 InputEncoder = util.ComposeErr(
	wt.AnyToF64,
	util.OkFunc(wa.EncodeF64),
)

var EncoderMap map[wa.ValueType]InputEncoder = map[wa.ValueType]InputEncoder{
	wa.ValueTypeI32: EncodeI32,
	wa.ValueTypeI64: EncodeI64,
	wa.ValueTypeF32: EncodeF32,
	wa.ValueTypeF64: EncodeF64,
}
