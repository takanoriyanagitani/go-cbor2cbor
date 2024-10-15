package cbor2primitive

import (
	"context"
	"log"
	"time"
)

type CborToAny func(context.Context, []byte) (any, error)

func (a CborToAny) Fallback(fallback CborToAny) CborToAny {
	switch a {
	case nil:
		return fallback
	default:
		return a
	}
}

type CborToPrimitive[T any] func([]byte) (T, error)

func (p CborToPrimitive[T]) ToConverter() CborToAny {
	return func(_ context.Context, byteString []byte) (any, error) {
		return p(byteString)
	}
}

var CborToAnyFallback CborToAny = func(
	_ context.Context,
	byteString []byte,
) (any, error) {
	return byteString, nil
}

type ConverterMapIx func(arrayIndex int) CborToAny

func (m ConverterMapIx) ConvertAll(
	ctx context.Context,
	arr [][]byte,
	fallback CborToAny,
) ([]any, error) {
	var ret []any = make([]any, 0, len(arr))
	var flbk CborToAny = fallback.Fallback(CborToAnyFallback)
	for i, val := range arr {
		select {
		case <-ctx.Done():
			return ret, ctx.Err()
		default:
			break
		}

		var cnv CborToAny = m(i).Fallback(flbk)
		if nil == val {
			ret = append(ret, nil)
			continue
		}
		converted, e := cnv(ctx, val)
		if nil != e {
			log.Printf("raw arr len: %v\n", len(arr))
			return ret, e
		}
		ret = append(ret, converted)
	}
	return ret, nil
}

func ConverterMapIxFromMap(m map[int]CborToAny) ConverterMapIx {
	return func(ix int) CborToAny {
		cta, ok := m[ix]
		switch ok {
		case true:
			return cta
		default:
			return CborToAnyFallback
		}
	}
}

var CborToAnyFromByteMap map[byte]CborToAny = map[byte]CborToAny{
	'b': CborToPrimitive[bool](CborToBool).ToConverter(),

	'H': CborToPrimitive[uint16](CborToUint16BE).ToConverter(),
	'I': CborToPrimitive[uint32](CborToUint32BE).ToConverter(),
	'Q': CborToPrimitive[uint64](CborToUint64BE).ToConverter(),

	'u': CborToPrimitive[[2]uint64](CborToUuidPair).ToConverter(),
	'U': CborToPrimitive[string](CborToUuidString).ToConverter(),

	't': CborToPrimitive[time.Time](CborToUnixtimeSecondsFloatBE).ToConverter(),

	'h': CborToPrimitive[int16](CborToInt16BE).ToConverter(),
	'i': CborToPrimitive[int32](CborToInt32BE).ToConverter(),
	'q': CborToPrimitive[int64](CborToInt64BE).ToConverter(),

	'd': CborToPrimitive[float64](CborToDoubleBE).ToConverter(),
	'f': CborToPrimitive[float32](CborToFloatBE).ToConverter(),

	's': CborToPrimitive[string](CborToStringUtf8).ToConverter(),
}

func CborToAnyFromByte(b byte) CborToAny {
	cnv, ok := CborToAnyFromByteMap[b]
	switch ok {
	case true:
		return cnv
	default:
		return CborToAnyFallback
	}
}

func ConverterMapIxFromChars(chars []uint8) ConverterMapIx {
	var m map[int]CborToAny = map[int]CborToAny{}
	for i, ch := range chars {
		m[i] = CborToAnyFromByte(ch)
	}
	return ConverterMapIxFromMap(m)
}

func ConverterMapIxFromFormatString(fstr string) ConverterMapIx {
	return ConverterMapIxFromChars([]byte(fstr))
}
