package main

import (
	"bufio"
	"context"
	"io"
	"iter"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	w0 "github.com/tetratelabs/wazero"

	cc "github.com/takanoriyanagitani/go-cbor2cbor"

	it "github.com/takanoriyanagitani/go-cbor2cbor/iter"

	cnv "github.com/takanoriyanagitani/go-cbor2cbor/conv"

	ia "github.com/takanoriyanagitani/go-cbor2cbor/iter/arr2cbor"
	aa "github.com/takanoriyanagitani/go-cbor2cbor/iter/arr2cbor/amacker"

	ci "github.com/takanoriyanagitani/go-cbor2cbor/iter/cbor2arr"
	ca "github.com/takanoriyanagitani/go-cbor2cbor/iter/cbor2arr/amacker"

	pw "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive/wasm"
	uw "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive/wasm/typed/unary/wazero"
)

type app struct {
	ci.CborToArrays
	ia.StreamWriter
	cnv.ConvMap
	DefaultConverter cnv.Converter
}

func (a app) ToArrayToCbor() ia.ArrayToCbor {
	return a.StreamWriter.ToArrayToCbor(
		a.ConvMap,
		a.DefaultConverter,
	)
}

func (a app) ToCborToMapd() cc.CborToMapd {
	return a.
		ToArrayToCbor().
		ToCborToMapd(a.CborToArrays)
}

func (a app) OutputAll(ctx context.Context) error {
	return a.ToCborToMapd()(ctx)
}

type InOut struct {
	io.Reader
	io.Writer
}

func (i InOut) ToCborToArrays() ci.CborToArrays {
	return ca.CborToArrNew(i.Reader).AsCborToArrays()
}

func (i InOut) ToStreamWriter() ia.StreamWriter {
	return aa.ArrToCborNew(i.Writer).ToStreamWriter()
}

type WasmConfig struct {
	Indices []uint32
	pw.WasmStore
	w0.Runtime
}

func (w WasmConfig) Close(ctx context.Context) error {
	return w.Runtime.Close(ctx)
}

func (w WasmConfig) ToModules(ctx context.Context) (uw.Modules, error) {
	return uw.ModulesNewDefaultFS(
		ctx,
		w.Indices,
		w.WasmStore,
		w.Runtime,
	)
}

func rdr2wtr(
	ctx context.Context,
	indices []uint32,
	wstore pw.WasmStore,
	rtm w0.Runtime,
	rdr io.Reader,
	wtr io.Writer,
) error {
	wcfg := WasmConfig{
		Indices:   indices,
		WasmStore: wstore,
		Runtime:   rtm,
	}
	defer wcfg.Close(ctx)

	modules, e := wcfg.ToModules(ctx)
	if nil != e {
		return e
	}
	defer modules.CloseAll(ctx)

	convmap, e := modules.ToConverterMap(ctx)
	if nil != e {
		return e
	}

	var cmap cnv.ConvMap = cnv.ConvMap(convmap)

	var convdefault cnv.Converter = uw.ConverterIdentity

	iout := InOut{
		Reader: rdr,
		Writer: wtr,
	}
	var c2a ci.CborToArrays = iout.ToCborToArrays()
	var sw ia.StreamWriter = iout.ToStreamWriter()

	ap := app{
		CborToArrays:     c2a,
		StreamWriter:     sw,
		ConvMap:          cmap,
		DefaultConverter: convdefault,
	}

	return ap.OutputAll(ctx)
}

func stdin2stdout(
	ctx context.Context,
	indices []uint32,
	wstore pw.WasmStore,
) error {
	var i io.Reader = os.Stdin
	var br io.Reader = bufio.NewReader(i)

	var o io.Writer = os.Stdout
	var bw *bufio.Writer = bufio.NewWriter(o)
	defer bw.Flush()

	var rtm w0.Runtime = w0.NewRuntime(ctx)
	defer rtm.Close(ctx)

	return rdr2wtr(
		ctx,
		indices,
		wstore,
		rtm,
		br,
		bw,
	)
}

func sub(
	ctx context.Context,
	indicesKey string,
	wasmdirKey string,
) error {
	var wstore pw.WasmStore = pw.WasmStoreNewFsTrustedFromEnv(wasmdirKey)

	var indices string = os.Getenv(indicesKey)
	var splited []string = strings.Split(indices, ",")
	var idxis iter.Seq[string] = slices.Values(splited)
	var idxip iter.Seq[it.Pair[error, uint32]] = it.Map(
		idxis,
		func(idxs string) it.Pair[error, uint32] {
			parsed, e := strconv.ParseUint(idxs, 10, 32)
			return it.Pair[error, uint32]{
				Left:  e,
				Right: uint32(parsed),
			}
		},
	)
	su, e := it.CollectErr(idxip)
	if nil != e {
		if "" != indices {
			return e
		}

		return stdin2stdout(ctx, nil, wstore)
	}

	return stdin2stdout(ctx, su, wstore)
}

func main() {
	e := sub(
		context.Background(),
		"ENV_COLUMN_INDICES",
		"ENV_WASM_MODULES_DIR",
	)
	if nil != e {
		log.Printf("%v\n", e)
	}
}
