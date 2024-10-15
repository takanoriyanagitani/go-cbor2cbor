package cbor2arr2cbor

import (
	"context"
	"errors"
	"io"

	a2c "github.com/takanoriyanagitani/go-cbor2cbor/arr2cbor"
	c2a "github.com/takanoriyanagitani/go-cbor2cbor/cbor2arr"
	cnv "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive"
)

type CborToArrayToMapdToCbor struct {
	// Parses the CBOR bytes and gets an array.
	c2a.CborToArray

	// Converts the input array and saves the converted array to the buffer.
	cnv.Converters

	// Serializes the array.
	a2c.ArrayToCbor
}

func (c CborToArrayToMapdToCbor) Convert(
	ctx context.Context,
	ibuf *[]any,
	obuf *[]any,
) error {
	ie := c.CborToArray(ibuf)
	if nil != ie {
		return ie
	}
	ce := c.Converters.ConvertAll(ctx, *ibuf, obuf)
	if nil != ce {
		return ce
	}
	return c.ArrayToCbor(*obuf)
}

func (c CborToArrayToMapdToCbor) ConvertAll(ctx context.Context) error {
	var ibuf []any
	var obuf []any
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ibuf = ibuf[:0]
		obuf = obuf[:0]
		e := c.Convert(ctx, &ibuf, &obuf)
		if nil != e {
			if !errors.Is(e, io.EOF) {
				return e
			}
			return nil
		}
	}
}
