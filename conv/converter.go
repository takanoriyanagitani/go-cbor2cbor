package conv

import (
	"context"
)

type Converter func(context.Context, any) (any, error)

type ConvMap map[uint32]Converter

func (m ConvMap) GetConverter(idx uint32, alt Converter) Converter {
	conv, found := m[idx]
	switch found {
	case true:
		return conv
	default:
		return alt
	}
}
