package any2cbor

import (
	"io"

	ac "github.com/fxamacker/cbor/v2"

	a2w "github.com/takanoriyanagitani/go-cbor2cbor/any2writer"
)

type AnyArrayToCbor struct {
	*ac.Encoder
}

func (a AnyArrayToCbor) Encode(arr []any) error {
	return a.Encoder.Encode(arr)
}

func AnyArrayToCborNew(wtr io.Writer) func([]any) error {
	em, err := ac.CanonicalEncOptions().EncMode()
	var a2c AnyArrayToCbor
	if nil == err {
		a2c.Encoder = em.NewEncoder(wtr)
	}
	return func(arr []any) error {
		if nil != err {
			return err
		}
		return a2c.Encode(arr)
	}
}

var AnyToWriterNew a2w.AnyToWriterNew = AnyArrayToCborNew
