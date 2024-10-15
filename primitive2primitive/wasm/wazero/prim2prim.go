package prim2prim

import (
	"context"
	"errors"
	"fmt"

	wz "github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"

	p2p "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive"

	w2p "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive/wasm"
)

var (
	ErrInvalidResult   error = errors.New("invalid result length")
	ErrInvalidFunction error = errors.New("invalid function")
	ErrInstantiate     error = errors.New("unable to create an instance")
	ErrWazero          error = errors.New("runtime error")
	ErrInvalidType     error = errors.New("invalid function type")
	ErrFunctionMissing error = errors.New("function not found")
)

type RawCaller struct {
	wa.Function
}

func (r RawCaller) ToConverter(typ any) p2p.PrimitiveToPrimitiveA {
	switch typ.(type) {
	case p2p.FloatToFloat32:
		var cnv p2p.FloatToFloat32 = r.ToFloatToFloat32()
		return cnv.ToAnyF()
	case p2p.IntToInt32:
		var cnv p2p.IntToInt32 = r.ToIntToInt32()
		return cnv.ToAnyI()
	case p2p.UintToUint32:
		var cnv p2p.UintToUint32 = r.ToUintToUint32()
		return cnv.ToAnyU()
	case p2p.IntToInt64:
		var cnv p2p.IntToInt64 = r.ToIntToInt64()
		return cnv.ToAnyI()
	case p2p.UintToUint64:
		var cnv p2p.UintToUint64 = r.ToUintToUint64()
		f2f := p2p.PrimitiveToPrimitive[uint64, uint64](cnv)
		return f2f.ToAny()
	case p2p.FloatToFloat64:
		var cnv p2p.FloatToFloat64 = r.ToFloatToFloat64()
		return cnv.ToAnyD()
	case p2p.LongToDouble:
		var cnv p2p.LongToDouble = r.ToLongToDouble()
		return cnv.ToAnyL()
	default:
		return func(_ context.Context, _ any) (any, error) {
			return 0, ErrInvalidType
		}
	}
}

func (r RawCaller) Single(ctx context.Context, i uint64) (uint64, error) {
	if nil == r.Function {
		return 0, ErrInvalidFunction
	}

	results, e := r.Function.Call(ctx, i)
	if nil != e {
		return 0, errors.Join(ErrWazero, e)
	}

	switch len(results) {
	case 1:
		return results[0], nil
	default:
		return 0, fmt.Errorf(
			"unexpected number of results(%v): %w",
			len(results),
			ErrInvalidResult,
		)
	}
	return 0, nil
}

func (r RawCaller) ToIntToInt32() p2p.IntToInt32 {
	return func(ctx context.Context, i int32) (int32, error) {
		var input uint64 = wa.EncodeI32(i)
		u, e := r.Single(ctx, input)
		if nil != e {
			return 0, e
		}
		return wa.DecodeI32(u), nil
	}
}

func (r RawCaller) ToUintToUint32() p2p.UintToUint32 {
	return func(ctx context.Context, i uint32) (uint32, error) {
		var input uint64 = wa.EncodeU32(i)
		u, e := r.Single(ctx, input)
		if nil != e {
			return 0, e
		}
		return wa.DecodeU32(u), nil
	}
}

func (r RawCaller) ToFloatToFloat32() p2p.FloatToFloat32 {
	return func(ctx context.Context, i float32) (float32, error) {
		var input uint64 = wa.EncodeF32(i)
		u, e := r.Single(ctx, input)
		if nil != e {
			return 0, e
		}
		return wa.DecodeF32(u), nil
	}
}

func (r RawCaller) ToIntToInt64() p2p.IntToInt64 {
	return func(ctx context.Context, i int64) (int64, error) {
		var input uint64 = wa.EncodeI64(i)
		u, e := r.Single(ctx, input)
		if nil != e {
			return 0, e
		}
		return int64(u), nil
	}
}

func (r RawCaller) ToUintToUint64() p2p.UintToUint64 {
	return func(ctx context.Context, i uint64) (uint64, error) {
		var input uint64 = i
		u, e := r.Single(ctx, input)
		if nil != e {
			return 0, e
		}
		return u, nil
	}
}

func (r RawCaller) ToFloatToFloat64() p2p.FloatToFloat64 {
	return func(ctx context.Context, i float64) (float64, error) {
		var input uint64 = wa.EncodeF64(i)
		u, e := r.Single(ctx, input)
		if nil != e {
			return 0, e
		}
		return wa.DecodeF64(u), nil
	}
}

func (r RawCaller) ToLongToDouble() p2p.LongToDouble {
	return func(ctx context.Context, i int64) (float64, error) {
		var input uint64 = wa.EncodeI64(i)
		u, e := r.Single(ctx, input)
		if nil != e {
			return 0, e
		}
		return wa.DecodeF64(u), nil
	}
}

type Instance struct {
	wz.CompiledModule
	wa.Module
	FunctionName string
	FunctionType any
}

func (i Instance) CloseC(ctx context.Context) error {
	if nil == i.CompiledModule {
		return nil
	}
	return i.CompiledModule.Close(ctx)
}

func (i Instance) CloseM(ctx context.Context) error {
	if nil == i.Module {
		return nil
	}
	return i.Module.Close(ctx)
}

func (i Instance) Close(ctx context.Context) error {
	em := i.CloseM(ctx)
	ec := i.CloseC(ctx)
	return errors.Join(em, ec)
}

func (i Instance) ToFunction(
	ctx context.Context,
) (p2p.PrimitiveToPrimitiveA, error) {
	if nil == i.Module {
		return nil, nil
	}
	var f wa.Function = i.Module.ExportedFunction(i.FunctionName)
	if nil == f {
		return nil, ErrFunctionMissing
	}
	raw := RawCaller{Function: f}
	return raw.ToConverter(i.FunctionType), nil
}

type Instances struct {
	wz.Runtime
	Instances []Instance
}

func (i Instances) CloseInstances(ctx context.Context) []error {
	var earr []error
	for _, instance := range i.Instances {
		earr = append(earr, instance.Close(ctx))
	}
	return earr
}

func (i Instances) Close(ctx context.Context) error {
	var errs []error = i.CloseInstances(ctx)
	var err error = i.Runtime.Close(ctx)
	return errors.Join(append(errs, err)...)
}

func (i Instances) ToFunctionArray(
	ctx context.Context,
) ([]p2p.PrimitiveToPrimitiveA, error) {
	var ret []p2p.PrimitiveToPrimitiveA
	for _, instance := range i.Instances {
		cnv, e := instance.ToFunction(ctx)
		if nil != e {
			return nil, e
		}
		ret = append(ret, cnv)
	}
	return ret, nil
}

func (i Instances) AsIndexToFunction() w2p.IndexToFunction {
	return i
}

func InstancesNewMust(
	ctx context.Context,
	config []w2p.IndexToFunc,
	store w2p.WasmStore,
	rtime wz.Runtime,
	mcfg wz.ModuleConfig,
) Instances {
	var imap map[uint8]Instance = map[uint8]Instance{}
	var imax uint8 = 0
	for _, cfg := range config {
		var ix uint8 = cfg.Index
		var moduleName string = cfg.Module
		var funcName string = cfg.Name
		var funcType any = cfg.ToFunctionType()

		imax = max(imax, ix)

		wasmBytes, e := store(ctx, moduleName)
		if nil != e {
			panic(e)
		}

		compiled, e := rtime.CompileModule(ctx, wasmBytes)
		if nil != e {
			panic(e)
		}

		instance, e := rtime.InstantiateModule(
			ctx,
			compiled,
			mcfg,
		)
		if nil != e {
			panic(e)
		}

		inst := Instance{
			CompiledModule: compiled,
			Module:         instance,
			FunctionName:   funcName,
			FunctionType:   funcType,
		}

		imap[ix] = inst
	}

	var instances []Instance = make([]Instance, imax+1)

	for i := range instances {
		instance, found := imap[uint8(i)]
		switch found {
		case true:
			instances[i] = instance
		default:
			instances[i] = Instance{}
		}
	}

	return Instances{
		Runtime:   rtime,
		Instances: instances,
	}
}

func InstancesNewDefaultMust(
	ctx context.Context,
	config []w2p.IndexToFunc,
	store w2p.WasmStore,
) Instances {
	var rtm wz.Runtime = wz.NewRuntime(ctx)
	var cfg wz.ModuleConfig = wz.NewModuleConfig()
	return InstancesNewMust(
		ctx,
		config,
		store,
		rtm,
		cfg,
	)
}
