package amacker

import (
	"io"

	ac "github.com/fxamacker/cbor/v2"

	r2b "github.com/takanoriyanagitani/go-cbor2cbor/reader2bytestrings"
)

type DecodeToByteStrings struct {
	*ac.Decoder
}

func (d DecodeToByteStrings) DecodeToByteStrings(bs *[][]byte) error {
	return d.Decoder.Decode(&bs)
}

func DecodeToByteStringsNew(rdr io.Reader) DecodeToByteStrings {
	var dec *ac.Decoder = ac.NewDecoder(rdr)
	return DecodeToByteStrings{dec}
}

func RdrToByteStrsNew() r2b.RdrToByteStrs {
	return func(rdr io.Reader) func(*[][]byte) error {
		var d2bs DecodeToByteStrings = DecodeToByteStringsNew(rdr)
		return d2bs.DecodeToByteStrings
	}
}
