package arr2cbor

import (
	"context"
	"iter"

	cc "github.com/takanoriyanagitani/go-cbor2cbor"

	cnv "github.com/takanoriyanagitani/go-cbor2cbor/conv"

	ci "github.com/takanoriyanagitani/go-cbor2cbor/iter/cbor2arr"
)

type ArrayToCbor func(context.Context, []any) error

func (a ArrayToCbor) OutputAll(ctx context.Context, i iter.Seq[[]any]) error {
	for arr := range i {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		e := a(ctx, arr)
		if nil != e {
			return e
		}
	}
	return nil
}

func (a ArrayToCbor) ToCborToMapd(c2a ci.CborToArrays) cc.CborToMapd {
	return func(ctx context.Context) error {
		var i iter.Seq[[]any] = c2a()
		return a.OutputAll(ctx, i)
	}
}

type StreamWriter struct {
	StartArray func(context.Context) error
	EndArray   func(context.Context) error
	AddElement func(context.Context, any) error
}

func (s StreamWriter) ToArrayToCbor(
	cmap cnv.ConvMap,
	alt cnv.Converter,
) ArrayToCbor {
	return func(ctx context.Context, arr []any) error {
		e := s.StartArray(ctx)
		if nil != e {
			return e
		}

		for i, item := range arr {
			var conv cnv.Converter = cmap.GetConverter(uint32(i), alt)
			converted, e := conv(ctx, item)
			if nil != e {
				return e
			}

			e = s.AddElement(ctx, converted)
			if nil != e {
				return e
			}
		}
		return s.EndArray(ctx)
	}
}
