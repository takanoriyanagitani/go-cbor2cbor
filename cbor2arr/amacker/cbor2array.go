package cbor2arr

import (
	"io"

	fa "github.com/fxamacker/cbor/v2"

	c2a "github.com/takanoriyanagitani/go-cbor2cbor/cbor2arr"
)

type CborToArr struct {
	*fa.Decoder
}

func (c CborToArr) ToArray(buf *[]any) error {
	return c.Decoder.Decode(buf)
}

func (c CborToArr) AsConverter() c2a.CborToArray {
	return c.ToArray
}

func CborToArrNew(rdr io.Reader) CborToArr {
	return CborToArr{
		Decoder: fa.NewDecoder(rdr),
	}
}
