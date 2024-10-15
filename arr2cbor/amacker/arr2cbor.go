package arr2cbor

import (
	"io"

	fa "github.com/fxamacker/cbor/v2"

	a2c "github.com/takanoriyanagitani/go-cbor2cbor/arr2cbor"
)

type ArrToCbor struct {
	*fa.Encoder
}

func (a ArrToCbor) Serialize(arr []any) error {
	return a.Encoder.Encode(arr)
}

func (a ArrToCbor) AsSerializer() a2c.ArrayToCbor {
	return a.Serialize
}

func ArrToCborNew(mode fa.EncMode) func(io.Writer) ArrToCbor {
	return func(wtr io.Writer) ArrToCbor {
		return ArrToCbor{
			Encoder: mode.NewEncoder(wtr),
		}
	}
}

func ArrToCborFromOpts(opts fa.EncOptions) func(io.Writer) (ArrToCbor, error) {
	mode, e := opts.EncMode()
	return func(wtr io.Writer) (ArrToCbor, error) {
		if nil != e {
			return ArrToCbor{}, e
		}
		return ArrToCborNew(mode)(wtr), nil
	}
}
