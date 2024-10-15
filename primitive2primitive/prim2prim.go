package prim2prim

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidArgument error = errors.New("invalid argument")
)

type PrimitiveToPrimitive[I, O any] func(context.Context, I) (O, error)

func (p PrimitiveToPrimitive[I, O]) ToAny() PrimitiveToPrimitiveA {
	return func(ctx context.Context, i any) (any, error) {
		switch input := i.(type) {
		case I:
			return p(ctx, input)
		default:
			return 0, fmt.Errorf(
				"unexpected type(%v): %w",
				input,
				ErrInvalidArgument,
			)
		}
	}
}

type PrimitiveToPrimitiveA func(context.Context, any) (any, error)

// The identity function.
var PrimitiveToPrimitiveAi PrimitiveToPrimitiveA = func(
	_ context.Context,
	i any,
) (any, error) {
	return i, nil
}

type IntToInt32 PrimitiveToPrimitive[int32, int32]
type UintToUint32 PrimitiveToPrimitive[uint32, uint32]
type FloatToFloat32 PrimitiveToPrimitive[float32, float32]

type IntToInt64 PrimitiveToPrimitive[int64, int64]
type UintToUint64 PrimitiveToPrimitive[uint64, uint64]
type FloatToFloat64 PrimitiveToPrimitive[float64, float64]

type LongToDouble PrimitiveToPrimitive[int64, float64]

type TimeFromUnixtime PrimitiveToPrimitive[float64, time.Time]

func (f2t TimeFromUnixtime) ToAnyT() PrimitiveToPrimitiveA {
	return func(ctx context.Context, i any) (any, error) {
		switch input := i.(type) {
		case float32:
			return f2t(ctx, float64(input))
		case float64:
			return f2t(ctx, input)
		case int64:
			return f2t(ctx, float64(input))
		case uint64:
			return f2t(ctx, float64(input))
		default:
			return 0, fmt.Errorf(
				"expected float32 or float64 or int64 or uint64. got=%v: %w",
				input,
				ErrInvalidArgument,
			)
		}
	}
}

func (i2i IntToInt32) ToAnyI() PrimitiveToPrimitiveA {
	return func(ctx context.Context, i any) (any, error) {
		switch input := i.(type) {
		case int32:
			return i2i(ctx, input)
		case int64:
			return i2i(ctx, int32(input))
		default:
			return 0, fmt.Errorf(
				"expected int32 or int64. got=%v: %w",
				input,
				ErrInvalidArgument,
			)
		}
	}
}

func (f2f FloatToFloat64) ToAnyD() PrimitiveToPrimitiveA {
	return func(ctx context.Context, i any) (any, error) {
		switch input := i.(type) {
		case int64:
			return f2f(ctx, float64(input))
		case uint64:
			return f2f(ctx, float64(input))
		case float32:
			return f2f(ctx, float64(input))
		default:
			return 0, fmt.Errorf(
				"expected int64 or uint64 or float32. got=%v: %w",
				input,
				ErrInvalidArgument,
			)
		}
	}
}

func (l2d LongToDouble) ToAnyL() PrimitiveToPrimitiveA {
	return func(ctx context.Context, i any) (any, error) {
		switch input := i.(type) {
		case int32:
			return l2d(ctx, int64(input))
		case int64:
			return l2d(ctx, input)
		case uint64:
			return l2d(ctx, int64(input))
		case float64:
			return l2d(ctx, int64(input))
		default:
			return 0, fmt.Errorf(
				"expected int32 or int64 or uint64. got=%v: %w",
				input,
				ErrInvalidArgument,
			)
		}
	}
}

func (i2i IntToInt64) ToAnyI() PrimitiveToPrimitiveA {
	return func(ctx context.Context, i any) (any, error) {
		switch input := i.(type) {
		case int32:
			return i2i(ctx, int64(input))
		case int64:
			return i2i(ctx, input)
		case uint64:
			return i2i(ctx, int64(input))
		default:
			return 0, fmt.Errorf(
				"expected int32 or int64 or uint64. got=%v: %w",
				input,
				ErrInvalidArgument,
			)
		}
	}
}

func (u2u UintToUint32) ToAnyU() PrimitiveToPrimitiveA {
	return func(ctx context.Context, i any) (any, error) {
		switch input := i.(type) {
		case uint32:
			return u2u(ctx, input)
		case uint64:
			return u2u(ctx, uint32(input))
		default:
			return 0, fmt.Errorf(
				"expected uint32 or uint64, got=%v: %w",
				input,
				ErrInvalidArgument,
			)
		}
	}
}

func (f2f FloatToFloat32) ToAnyF() PrimitiveToPrimitiveA {
	return func(ctx context.Context, i any) (any, error) {
		switch input := i.(type) {
		case float32:
			return f2f(ctx, input)
		case float64:
			return f2f(ctx, float32(input))
		default:
			return 0, fmt.Errorf(
				"expected float32 or float64, got=%v: %w",
				input,
				ErrInvalidArgument,
			)
		}
	}
}

type Converters []PrimitiveToPrimitiveA

func (c Converters) GetConverter(i int) PrimitiveToPrimitiveA {
	var sz int = len(c)
	if i < sz {
		return c[i]
	}
	return PrimitiveToPrimitiveAi
}

func (c Converters) ConvertAll(
	ctx context.Context,
	inputs []any,
	buf *[]any,
) error {
	for i, input := range inputs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var cnv PrimitiveToPrimitiveA = c.GetConverter(i)
		if nil == cnv {
			var output any = input
			*buf = append(*buf, output)
			continue
		}

		if nil == input {
			var output any = nil
			*buf = append(*buf, output)
			continue
		}

		output, e := cnv(ctx, input)
		if nil != e {
			return e
		}
		*buf = append(*buf, output)
	}
	return nil
}
