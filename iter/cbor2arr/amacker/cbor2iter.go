package cbor2iter

import (
	"io"
	"iter"

	ac "github.com/fxamacker/cbor/v2"

	ic "github.com/takanoriyanagitani/go-cbor2cbor/iter/cbor2arr"
)

type CborToArr struct {
	*ac.Decoder
}

func CborToArrNew(rdr io.Reader) CborToArr {
	return CborToArr{Decoder: ac.NewDecoder(rdr)}
}

func (c CborToArr) ToArrays() iter.Seq[[]any] {
	return func(yield func([]any) bool) {
		var buf []any
		var err error
		for {
			clear(buf)
			buf = buf[:0]
			err = c.Decoder.Decode(&buf)
			if nil != err {
				return
			}

			if !yield(buf) {
				return
			}
		}
	}
}

func (c CborToArr) AsCborToArrays() ic.CborToArrays {
	return c.ToArrays
}
