package cbor2primitive

import (
	"errors"

	"unicode/utf8"
)

var (
	ErrInvalidStringUtf8 error = errors.New("invalid utf8 string")
)

func CborToStringUtf8(byteString []byte) (string, error) {
	var valid bool = utf8.Valid(byteString)
	switch valid {
	case true:
		return string(byteString), nil
	default:
		return "", ErrInvalidStringUtf8
	}
}
