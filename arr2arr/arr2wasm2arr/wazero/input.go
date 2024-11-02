package arr2wz2arr

import (
	"context"

	wa "github.com/tetratelabs/wazero/api"
)

type SetInputSize struct {
	wa.Function
}

func (s SetInputSize) SetInput(
	ctx context.Context,
	sz uint32,
	init uint8,
) (uint32, error) {
	var encoded uint64 = wa.EncodeU32(sz)
	results, e := s.Function.Call(ctx, encoded, wa.EncodeU32(uint32(init)))

	result := Result{
		Err: e,
		Raw: results,
	}
	i, e := result.GetInt32i(ctx)
	return IntToUintNonNeg(i, e)
}

func (s SetInputSize) SetInputDefault(
	ctx context.Context,
	sz uint32,
) (uint32, error) {
	return s.SetInput(ctx, sz, 0)
}

type InputOffset struct {
	wa.Function
}

func (o InputOffset) GetInputOffset(ctx context.Context) (uint32, error) {
	results, e := o.Function.Call(ctx)
	result := Result{
		Err: e,
		Raw: results,
	}
	i, e := result.GetInt32i(ctx)
	return IntToUintNonNeg(i, e)
}
