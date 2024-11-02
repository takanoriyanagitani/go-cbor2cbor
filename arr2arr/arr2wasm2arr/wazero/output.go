package arr2wz2arr

import (
	"context"

	wa "github.com/tetratelabs/wazero/api"
)

type EstimateOutputSize struct {
	wa.Function
}

func (o EstimateOutputSize) Estimate(ctx context.Context) (uint32, error) {
	results, e := o.Function.Call(ctx)
	result := Result{
		Err: e,
		Raw: results,
	}
	i, e := result.GetInt32i(ctx)
	return IntToUintNonNeg(i, e)
}

type SetOutputSize struct {
	wa.Function
}

func (o SetOutputSize) SetOutSize(ctx context.Context, sz uint32) (uint32, error) {
	var encoded uint64 = wa.EncodeU32(sz)
	results, e := o.Function.Call(ctx, encoded)
	result := Result{
		Err: e,
		Raw: results,
	}
	i, e := result.GetInt32i(ctx)
	return IntToUintNonNeg(i, e)
}

type OutputOffset struct {
	wa.Function
}

func (o OutputOffset) GetOutputOffset(ctx context.Context) (uint32, error) {
	results, e := o.Function.Call(ctx)
	result := Result{
		Err: e,
		Raw: results,
	}
	i, e := result.GetInt32i(ctx)
	return IntToUintNonNeg(i, e)
}
