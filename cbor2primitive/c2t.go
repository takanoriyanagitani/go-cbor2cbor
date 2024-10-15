package cbor2primitive

import (
	"errors"
	"time"

	"encoding/binary"
	"math"
)

var (
	ErrInvalidUnixtimeLength error = errors.New("invalid unixtime")
)

func CborToUnixtimeSecondsFloatBE(byteString []byte) (time.Time, error) {
	var sz int = len(byteString)
	switch sz {
	case 8:
		var u uint64 = binary.BigEndian.Uint64(byteString)
		var f float64 = math.Float64frombits(u)
		var micros float64 = f * 1e6
		return time.UnixMicro(int64(micros)), nil
	default:
		return time.UnixMicro(0), ErrInvalidUnixtimeLength
	}
}
