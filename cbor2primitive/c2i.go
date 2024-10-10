package cbor2primitive

import (
	"errors"
	"log"

	"encoding/binary"
)

var (
	ErrInvalidInt16Length error = errors.New("invalid int16")
	ErrInvalidInt32Length error = errors.New("invalid int32")
	ErrInvalidInt64Length error = errors.New("invalid int64")
)

func CborToInt64BE(byteString []byte) (int64, error) {
	var sz int = len(byteString)
	switch sz {
	case 8:
		return int64(binary.BigEndian.Uint64(byteString)), nil
	default:
		return 0, ErrInvalidInt64Length
	}
}

func CborToInt32BE(byteString []byte) (int32, error) {
	var sz int = len(byteString)
	switch sz {
	case 4:
		return int32(binary.BigEndian.Uint32(byteString)), nil
	default:
		return 0, ErrInvalidInt32Length
	}
}

func CborToInt16BE(byteString []byte) (int16, error) {
	var sz int = len(byteString)
	switch sz {
	case 2:
		return int16(binary.BigEndian.Uint16(byteString)), nil
	default:
		log.Printf("byte string: %v\n", byteString)
		return 0, ErrInvalidInt16Length
	}
}
