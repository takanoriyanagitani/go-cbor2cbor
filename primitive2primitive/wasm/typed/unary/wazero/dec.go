package unary0

import (
	wa "github.com/tetratelabs/wazero/api"
)

type OutputDecoder func(uint64) any

func DecodeI32(result uint64) any { return wa.DecodeI32(result) }
func DecodeI64(result uint64) any { return int64(result) }
func DecodeF32(result uint64) any { return wa.DecodeF32(result) }
func DecodeF64(result uint64) any { return wa.DecodeF64(result) }

var DecoderMap map[wa.ValueType]OutputDecoder = map[wa.ValueType]OutputDecoder{
	wa.ValueTypeI32: DecodeI32,
	wa.ValueTypeI64: DecodeI64,
	wa.ValueTypeF32: DecodeF32,
	wa.ValueTypeF64: DecodeF64,
}
