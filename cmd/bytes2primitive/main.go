package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os"

	r2p "github.com/takanoriyanagitani/go-cbor2cbor/reader2primitives"

	a2w "github.com/takanoriyanagitani/go-cbor2cbor/any2writer"
	a2wa "github.com/takanoriyanagitani/go-cbor2cbor/any2writer/amacker"

	c2p "github.com/takanoriyanagitani/go-cbor2cbor/cbor2primitive"

	r2b "github.com/takanoriyanagitani/go-cbor2cbor/reader2bytestrings"
	r2ba "github.com/takanoriyanagitani/go-cbor2cbor/reader2bytestrings/amacker"
)

var cvtFormat string = os.Getenv("ENV_CVT_FMT")

var cnvrtMapI c2p.ConverterMapIx = c2p.ConverterMapIxFromFormatString(cvtFormat)

var rdr2bstrs r2b.RdrToByteStrs = r2ba.RdrToByteStrsNew()
var cnvFllbck c2p.CborToAny = nil // use default

var rdr2p r2p.ReaderToPrimitives = r2p.ReaderToPrimitives{
	RdrToByteStrs:  rdr2bstrs,
	ConverterMapIx: cnvrtMapI,
	Fallback:       cnvFllbck,
}

var any2wtr a2w.AnyToWriterNew = a2wa.AnyToWriterNew

type app struct {
	io.Reader
	io.Writer

	// io.Reader -> []any
	r2p.RdrToPrimitives

	// []any -> io.Writer
	a2w.AnyToWriterNew
}

func (a app) WithAnyToWtrNew(a2wn a2w.AnyToWriterNew) app {
	a.AnyToWriterNew = a2wn
	return a
}

func (a app) WithReaderToPrimitives(rdr2prim r2p.ReaderToPrimitives) app {
	a.RdrToPrimitives = rdr2prim.ToConverter()
	return a
}

func (a app) WithReader(r io.Reader) app {
	a.Reader = r
	return a
}

func (a app) WithWriter(w io.Writer) app {
	a.Writer = w
	return a
}

func (a app) ToSource() func(context.Context) ([]any, error) {
	return a.RdrToPrimitives(a.Reader)
}

func (a app) ToSink() func([]any) error {
	return a.AnyToWriterNew(a.Writer)
}

func (a app) CborByteStringsToPrimitives(ctx context.Context) error {
	var src func(context.Context) ([]any, error) = a.ToSource()
	var dst func([]any) error = a.ToSink()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		row, e := src(ctx)
		if nil != e {
			if errors.Is(e, io.EOF) {
				return nil
			}
			return e
		}

		e = dst(row)
		if nil != e {
			return e
		}
	}
}

var appDefault app = app{}.
	WithAnyToWtrNew(any2wtr).
	WithReaderToPrimitives(rdr2p)

func main() {
	var i io.Reader = os.Stdin
	var br io.Reader = bufio.NewReader(i)

	var o io.Writer = os.Stdout
	var bw *bufio.Writer = bufio.NewWriter(o)
	defer bw.Flush()

	var a app = appDefault.
		WithReader(br).
		WithWriter(bw)

	e := a.CborByteStringsToPrimitives(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
		panic(e)
	}
}
