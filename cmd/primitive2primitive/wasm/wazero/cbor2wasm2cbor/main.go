package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	fa "github.com/fxamacker/cbor/v2"

	a2c "github.com/takanoriyanagitani/go-cbor2cbor/arr2cbor"
	a2ca "github.com/takanoriyanagitani/go-cbor2cbor/arr2cbor/amacker"

	c2c "github.com/takanoriyanagitani/go-cbor2cbor/cbor2arr2mapd2cbor"

	c2a "github.com/takanoriyanagitani/go-cbor2cbor/cbor2arr"
	c2aa "github.com/takanoriyanagitani/go-cbor2cbor/cbor2arr/amacker"

	p2p "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive"
	w2p "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive/wasm"
	z2p "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive/wasm/wazero"
)

type app struct {
	w2p.IndexToFunction
	c2a.CborToArray
	a2c.ArrayToCbor
}

func (a app) ConvertAll(ctx context.Context) error {
	defer a.IndexToFunction.Close(ctx)
	functions, e := a.IndexToFunction.ToFunctionArray(ctx)
	if nil != e {
		return e
	}
	var cnv p2p.Converters = functions
	mapper := c2c.CborToArrayToMapdToCbor{
		CborToArray: a.CborToArray,
		Converters:  cnv,
		ArrayToCbor: a.ArrayToCbor,
	}
	return mapper.ConvertAll(ctx)
}

func appNewMust(
	ctx context.Context,
	cfg []w2p.IndexToFunc,
	store w2p.WasmStore,
	parse c2a.CborToArray,
	ser a2c.ArrayToCbor,
) app {
	var instances z2p.Instances = z2p.InstancesNewDefaultMust(
		ctx,
		cfg,
		store,
	)
	var i2f w2p.IndexToFunction = instances.AsIndexToFunction()
	return app{
		IndexToFunction: i2f,
		CborToArray:     parse,
		ArrayToCbor:     ser,
	}
}

func rdr2wtrMust(
	ctx context.Context,
	cfgStore w2p.IndexToFuncStore,
	store w2p.WasmStore,
	rdr io.Reader,
	wtr io.Writer,
) error {
	var parse c2a.CborToArray = c2aa.CborToArrNew(rdr).AsConverter()

	var opts fa.EncOptions = fa.CanonicalEncOptions()
	ser, e := a2ca.ArrToCborFromOpts(opts)(wtr)
	if nil != e {
		return e
	}

	cfg, e := cfgStore(ctx)
	if nil != e {
		return e
	}

	var cnv a2c.ArrayToCbor = ser.AsSerializer()
	var a app = appNewMust(
		ctx,
		cfg,
		store,
		parse,
		cnv,
	)
	return a.ConvertAll(ctx)
}

func stdin2stdout(ctx context.Context) error {
	var i io.Reader = os.Stdin
	var o io.Writer = os.Stdout

	var br io.Reader = bufio.NewReader(i)
	var bw *bufio.Writer = bufio.NewWriter(o)
	defer bw.Flush()

	var wstore w2p.WasmStore = w2p.WasmStoreNewFsTrustedFromEnv(
		"ENV_TRUSTED_WASM_DIR",
	)

	var rstore w2p.IndexToFuncStoreRaw = w2p.IndexToFuncStoreRawFsNewFromEnv(
		"ENV_TRUSTED_CFG_DIR_NAME",
		"ENV_TRUSTED_CFG_FILENAME",
	)

	var cstore w2p.IndexToFuncStore = rstore.ToJsonStore()

	return rdr2wtrMust(
		ctx,
		cstore,
		wstore,
		br,
		bw,
	)
}

func main() {
	e := stdin2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
