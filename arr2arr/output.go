package arr2arr

import (
	"context"
)

// The output array converted by the converter.
type OutputArray []any

// The serialized output(or the source of the [OutputArray]).
type OutputCbor []byte

// Writes the output CBOR bytes.
type OutputWriter func(context.Context, OutputCbor) error

// Checks the converted CBOR bytes
type OutputChecker func(context.Context, OutputCbor) error

func (c OutputChecker) ToCheckedWriter(wtr OutputWriter) OutputWriter {
	return func(ctx context.Context, unchecked OutputCbor) error {
		var err error = c(ctx, unchecked)
		if nil != err {
			return err
		}
		return wtr(ctx, unchecked)
	}
}

// Trusts the converted CBOR bytes(assumes the bytes is CBOR []any).
func OutputCheckerNop(_ context.Context, _ OutputCbor) error { return nil }
