package rdr2primitive

import (
	"context"
	"io"

	c2p "github.com/takanoriyanagitani/go-cbor2cbor/cbor2primitive"
	r2b "github.com/takanoriyanagitani/go-cbor2cbor/reader2bytestrings"
)

type RdrToPrimitives func(io.Reader) func(context.Context) ([]any, error)

type ReaderToPrimitives struct {
	r2b.RdrToByteStrs
	c2p.ConverterMapIx
	Fallback c2p.CborToAny
}

func (p ReaderToPrimitives) ToConverter() RdrToPrimitives {
	return func(rdr io.Reader) func(context.Context) ([]any, error) {
		var rdr2bstr func(*[][]byte) error = p.RdrToByteStrs(rdr)
		var buf [][]byte
		return func(ctx context.Context) ([]any, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			e := rdr2bstr(&buf)
			if nil != e {
				return nil, e
			}
			return p.ConverterMapIx.ConvertAll(
				ctx,
				buf,
				p.Fallback,
			)
		}
	}
}
