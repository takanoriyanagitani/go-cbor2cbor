package arr2arr

import (
	"context"
)

type ArrayToArray func(context.Context, InputArray) (OutputArray, error)

// Directly pass the input CBOR and gets the converted CBOR.
type RawArrayToRawArray func(context.Context, InputCbor) (OutputCbor, error)

func (r RawArrayToRawArray) ToArrayToRawArray(i2c InputToCbor) ArrayToRawArray {
	return func(ctx context.Context, i InputArray) (OutputCbor, error) {
		ser, e := i2c(ctx, i)
		if nil != e {
			return nil, e
		}

		return r(ctx, ser)
	}
}

// Pass the parsed CBOR array and gets the converted CBOR bytes(array).
type ArrayToRawArray func(context.Context, InputArray) (OutputCbor, error)
