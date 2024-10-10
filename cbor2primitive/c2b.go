package cbor2primitive

import (
	"errors"
)

var (
	ErrInvalidBoolLength error = errors.New("invalid bool")
)

func CborToBool(byteString []byte) (bool, error) {
	var sz int = len(byteString)
	switch sz {
	default:
		for _, b := range byteString {
			if 0 != b {
				return true, nil
			}
		}
		return false, nil
	}
}
