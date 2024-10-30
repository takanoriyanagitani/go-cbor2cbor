package typed

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInputType error = errors.New("invalid input type")
)

type SignedInteger interface {
	~int32 | ~int64
}

type Unsigned interface {
	~uint32 | ~uint64
}

type Integer interface {
	SignedInteger | Unsigned
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Integer | Float
}

func AnyToNumber[T Number](
	input any,
	i2t32 func(int32) T,
	i2t64 func(int64) T,
	u2t32 func(uint32) T,
	u2t64 func(uint64) T,
	f2t32 func(float32) T,
	f2t64 func(float64) T,
) (T, error) {
	switch i := input.(type) {
	case int32:
		return i2t32(i), nil
	case int64:
		return i2t64(i), nil
	case uint32:
		return u2t32(i), nil
	case uint64:
		return u2t64(i), nil
	case float32:
		return f2t32(i), nil
	case float64:
		return f2t64(i), nil
	default:
		return 0, fmt.Errorf("%w: rejected value=%v", ErrInvalidInputType, i)
	}
}
