package amacker_test

import (
	"testing"

	"bytes"
	"errors"
	"io"

	c2b "github.com/takanoriyanagitani/go-cbor2cbor/reader2bytestrings/amacker"
)

func TestR2b(t *testing.T) {
	t.Parallel()

	t.Run("DecodeToByteStrings", func(t *testing.T) {
		t.Parallel()

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			var d2b c2b.DecodeToByteStrings = c2b.DecodeToByteStringsNew(
				bytes.NewReader(nil),
			)

			var buf [][]byte

			e := d2b.DecodeToByteStrings(&buf)
			if !errors.Is(e, io.EOF) {
				t.Fatalf("unexpected err: %v\n", e)
			}
		})

		t.Run("i32,i32(null)", func(t *testing.T) {
			t.Parallel()

			var d2b c2b.DecodeToByteStrings = c2b.DecodeToByteStringsNew(
				bytes.NewReader([]byte{
					0x82, 0x44,
					0x00, 0x00,
					0x00, 0x01,
					0xf6,
				}),
			)

			var buf [][]byte

			e := d2b.DecodeToByteStrings(&buf)
			if nil != e {
				t.Fatalf("unexpected err: %v\n", e)
			}

			switch len(buf) {
			case 2:
				break
			default:
				t.Fatalf("unexpected col count: %v\n", len(buf))
			}

			var nb []byte = buf[1]
			if 0 != len(nb) {
				t.Fatalf("unexpected size: %v\n", len(nb))
			}

			var ib []byte = buf[0]
			if 4 != len(ib) {
				t.Fatalf("unexpected size: %v\n", len(ib))
			}
		})
	})
}
