package cbor2primitive

import (
	"errors"

	"encoding/binary"
	"math"
)

var (
	ErrInvalidDoubleLength error = errors.New("invalid double")
)

func CborToDoubleBE(byteString []byte) (float64, error) {
	var sz int = len(byteString)
	switch sz {
	case 8:
		var u uint64 = binary.BigEndian.Uint64(byteString)
		var f float64 = math.Float64frombits(u)
		return f, nil
	default:
		return 0.0, ErrInvalidDoubleLength
	}
}
