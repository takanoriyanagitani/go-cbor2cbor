package arr2wz2arr

import (
	"context"
	"errors"
	"fmt"

	wa "github.com/tetratelabs/wazero/api"

	util "github.com/takanoriyanagitani/go-cbor2cbor/util"
)

var (
	ErrInvalidResultType error = errors.New("invalid result type")
)

type Result struct {
	Err error
	Raw []uint64
}

func (r *Result) GetErr() error { return r.Err }

func (r *Result) GetRawResult(_ context.Context) (uint64, error) {
	switch len(r.Raw) {
	case 1:
		return r.Raw[0], r.GetErr()
	default:
		return 0, fmt.Errorf(
			"%w: wrong result length(%v): %v",
			ErrInvalidResultType,
			len(r.Raw),
			r.GetErr(),
		)
	}
}

func (r *Result) GetInt32u(ctx context.Context) (uint32, error) {
	return util.ComposeIO[uint64, uint32](
		r.GetRawResult,
		wa.DecodeU32,
	)(ctx)
}

func (r *Result) GetInt32i(ctx context.Context) (int32, error) {
	return util.ComposeIO[uint64, int32](
		r.GetRawResult,
		wa.DecodeI32,
	)(ctx)
}
