package arr2cbor

import (
	"context"
	"io"

	ac "github.com/fxamacker/cbor/v2"

	ia "github.com/takanoriyanagitani/go-cbor2cbor/iter/arr2cbor"
)

type ArrToCbor struct {
	*ac.Encoder
}

func ArrToCborNew(wtr io.Writer) ArrToCbor {
	return ArrToCbor{Encoder: ac.NewEncoder(wtr)}
}

func (a ArrToCbor) StartArray(_ context.Context) error {
	return a.Encoder.StartIndefiniteArray()
}

func (a ArrToCbor) EndArray(_ context.Context) error {
	return a.Encoder.EndIndefinite()
}

func (a ArrToCbor) AddElement(_ context.Context, item any) error {
	return a.Encoder.Encode(item)
}

func (a ArrToCbor) ToStreamWriter() ia.StreamWriter {
	return ia.StreamWriter{
		StartArray: a.StartArray,
		EndArray:   a.EndArray,
		AddElement: a.AddElement,
	}
}
