package cbor2primitive

import (
	"errors"

	"encoding/binary"
)

var (
	ErrInvalidUint16Length error = errors.New("invalid uint16")
	ErrInvalidUint32Length error = errors.New("invalid uint32")
	ErrInvalidUint64Length error = errors.New("invalid uint64")
)

func CborToUint64BE(byteString []byte) (uint64, error) {
	var sz int = len(byteString)
	switch sz {
	case 8:
		return binary.BigEndian.Uint64(byteString), nil
	default:
		return 0, ErrInvalidUint64Length
	}
}

func CborToUint32BE(byteString []byte) (uint32, error) {
	var sz int = len(byteString)
	switch sz {
	case 4:
		return binary.BigEndian.Uint32(byteString), nil
	default:
		return 0, ErrInvalidUint32Length
	}
}

func CborToUint16BE(byteString []byte) (uint16, error) {
	var sz int = len(byteString)
	switch sz {
	case 2:
		return binary.BigEndian.Uint16(byteString), nil
	default:
		return 0, ErrInvalidUint16Length
	}
}
