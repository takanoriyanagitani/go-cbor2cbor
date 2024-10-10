package cbor2primitive

import (
	"errors"

	"encoding/binary"
	"math"
)

var (
	ErrInvalidFloatLength error = errors.New("invalid float")
)

func CborToFloatBE(byteString []byte) (float32, error) {
	var sz int = len(byteString)
	switch sz {
	case 4:
		var u uint32 = binary.BigEndian.Uint32(byteString)
		var f float32 = math.Float32frombits(u)
		return f, nil
	default:
		return 0.0, ErrInvalidFloatLength
	}
}
