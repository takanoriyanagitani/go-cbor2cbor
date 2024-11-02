package arr2arr

import (
	"bytes"
	"context"
	"iter"
)

// The input array parsed from CBOR bytes.
type InputArray []any

// The serialized input(or the source of the [InputArray]).
type InputCbor []byte

// Parses the CBOR bytes and gets []any.
type InputToArray func(context.Context, InputCbor) ([]any, error)

// Serializes the input array and get serialized CBOR bytes(array).
type InputToCbor func(context.Context, InputArray) (InputCbor, error)

type InputToCborBuf func(context.Context, InputArray, *bytes.Buffer) error

func (b InputToCborBuf) ToInputToCbor() InputToCbor {
	var buf bytes.Buffer
	var err error
	return func(ctx context.Context, i InputArray) (InputCbor, error) {
		buf.Reset()
		err = b(ctx, i, &buf)
		return buf.Bytes(), err
	}
}

type InputArrayIter iter.Seq[[]any]

func (i InputArrayIter) OutputAll(
	ctx context.Context,
	a2r ArrayToRawArray,
	wtr OutputWriter,
) error {
	for input := range i {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		converted, e := a2r(ctx, input)
		if nil != e {
			return e
		}

		e = wtr(ctx, converted)
		if nil != e {
			return e
		}
	}
	return nil
}
