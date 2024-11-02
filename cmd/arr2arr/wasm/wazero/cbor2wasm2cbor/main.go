package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"iter"
	"log"
	"os"
	"strconv"

	w0 "github.com/tetratelabs/wazero"

	util "github.com/takanoriyanagitani/go-cbor2cbor/util"

	aa "github.com/takanoriyanagitani/go-cbor2cbor/arr2arr"
	wa "github.com/takanoriyanagitani/go-cbor2cbor/arr2arr/arr2wasm2arr"
	af "github.com/takanoriyanagitani/go-cbor2cbor/arr2arr/arr2wasm2arr/fs"
	aw "github.com/takanoriyanagitani/go-cbor2cbor/arr2arr/arr2wasm2arr/wazero"

	ca "github.com/takanoriyanagitani/go-cbor2cbor/iter/cbor2arr/amacker"

	ac "github.com/takanoriyanagitani/go-cbor2cbor/arr2cbor"
	am "github.com/takanoriyanagitani/go-cbor2cbor/arr2cbor/amacker"
)

func envkey2stringNew(key string) func(context.Context) (string, error) {
	return func(_ context.Context) (string, error) {
		return os.Getenv(key), nil
	}
}

func string2uint32altNew(alt uint32) func(string) uint32 {
	return func(s string) uint32 {
		u, e := strconv.ParseUint(s, 10, 32)
		switch e {
		case nil:
			return uint32(u)
		default:
			return alt
		}
	}
}

var envkey2wasmMaxBytes func(context.Context) (uint32, error) = util.ComposeIO(
	util.InputIO[string](envkey2stringNew("ENV_WASM_BYTES_MAX")),
	string2uint32altNew(af.WasmBytesMaxDefault),
)

type WasmSourceConfig struct {
	WasmDirname   string
	WasmName      string
	WasmSizeLimit uint32
}

func (w WasmSourceConfig) ToFsSource() af.FsSource {
	return af.FsSource{
		FS:       os.DirFS(w.WasmDirname),
		Basename: w.WasmName,
		MaxBytes: w.WasmSizeLimit,
	}
}

func (w WasmSourceConfig) ToWasmSource() wa.WasmSource {
	return w.ToFsSource().ToWasmSource()
}

type IoConfig struct {
	io.Reader
	io.Writer
}

func (c IoConfig) ToCborArrayIter() iter.Seq[[]any] {
	var br io.Reader = bufio.NewReader(c.Reader)
	return ca.CborToArrNew(br).ToArrays()
}

func (c IoConfig) ToInputArrIter() aa.InputArrayIter {
	var i iter.Seq[[]any] = c.ToCborArrayIter()
	return aa.InputArrayIter(i)
}

func (c IoConfig) ToBufWriter() *bufio.Writer {
	return bufio.NewWriter(c.Writer)
}

type App struct {
	aa.InputArrayIter
	aa.ArrayToRawArray
	aa.OutputWriter
}

func (a App) OutputAll(ctx context.Context) error {
	return a.InputArrayIter.OutputAll(
		ctx,
		a.ArrayToRawArray,
		a.OutputWriter,
	)
}

func rdr2wtr(ctx context.Context, icfg IoConfig, wcfg WasmSourceConfig) error {
	var wsrc wa.WasmSource = wcfg.ToWasmSource()

	var rtm aw.Runtime = aw.RuntimeNewDefault(w0.NewRuntime(ctx))

	compiled, e := rtm.WasmSourceIntoCompiled(
		ctx,
		wsrc,
	)
	if nil != e {
		return errors.Join(e, rtm.Close(ctx))
	}
	defer compiled.Close(ctx)

	conv, e := compiled.Instantiate(ctx)
	if nil != e {
		return e
	}
	defer conv.Close(ctx)

	var a2c am.AnyToCborToBuf = am.AnyToCborToBufDefault
	var a2b ac.ArrayToCborToBuffer = a2c.ToArrayToCborToBuffer()

	var r2r aa.RawArrayToRawArray = conv.ToRawArrToRawArr()
	var i2b aa.InputToCborBuf = func(
		ctx context.Context,
		ia aa.InputArray,
		buf *bytes.Buffer,
	) error {
		return a2b(ctx, ia, buf)
	}
	var i2c aa.InputToCbor = i2b.ToInputToCbor()
	var a2r aa.ArrayToRawArray = r2r.ToArrayToRawArray(i2c)

	var iai aa.InputArrayIter = icfg.ToInputArrIter()
	var wtr *bufio.Writer = icfg.ToBufWriter()
	defer wtr.Flush()

	var ochk aa.OutputChecker = aa.OutputCheckerNop
	var ow aa.OutputWriter = func(_ context.Context, o aa.OutputCbor) error {
		_, e := wtr.Write(o)
		return e
	}
	ow = ochk.ToCheckedWriter(ow)

	app := App{
		InputArrayIter:  iai,
		ArrayToRawArray: a2r,
		OutputWriter:    ow,
	}

	return app.OutputAll(ctx)
}

func stdin2stdout(ctx context.Context, wcfg WasmSourceConfig) error {
	icfg := IoConfig{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
	return rdr2wtr(ctx, icfg, wcfg)
}

func sub(ctx context.Context) error {
	wasmMaxBytes, e := envkey2wasmMaxBytes(ctx)
	if nil != e {
		return e
	}

	var wasmDir string = os.Getenv("ENV_WASM_MODULES_DIR")
	var wasmName string = os.Getenv("ENV_WASM_FILENAME")

	wcfg := WasmSourceConfig{
		WasmDirname:   wasmDir,
		WasmName:      wasmName,
		WasmSizeLimit: wasmMaxBytes,
	}
	return stdin2stdout(ctx, wcfg)
}

func main() {
	e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
