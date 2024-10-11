package cbor2primitive

import (
	"errors"

	"encoding/binary"
)

var (
	ErrInvalidUuidLength error = errors.New("invalid uuid")
)

func CborToUuidPair(byteString []byte) ([2]uint64, error) {
	var sz int = len(byteString)
	switch sz {
	case 16:
		var hi []byte = byteString[:8]
		var lo []byte = byteString[8:]
		var h uint64 = binary.BigEndian.Uint64(hi)
		var l uint64 = binary.BigEndian.Uint64(lo)
		return [2]uint64{h, l}, nil
	default:
		return [2]uint64{0, 0}, ErrInvalidUuidLength
	}
}
