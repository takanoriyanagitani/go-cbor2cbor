package fs2wasm

import (
	"bufio"
	"context"
	"io"
	"io/fs"

	aa "github.com/takanoriyanagitani/go-cbor2cbor/arr2arr/arr2wasm2arr"
)

const (
	WasmBytesMaxDefault uint32 = 16777216
)

type FsSource struct {
	fs.FS
	Basename string
	MaxBytes uint32
}

func (f FsSource) ReaderToBytes(r io.Reader) ([]byte, error) {
	limited := &io.LimitedReader{
		R: r,
		N: int64(f.MaxBytes),
	}
	return io.ReadAll(limited)
}

func (f FsSource) ToBytes() ([]byte, error) {
	file, e := f.FS.Open(f.Basename)
	if nil != e {
		return nil, e
	}
	defer file.Close()

	var br io.Reader = bufio.NewReader(file)
	return f.ReaderToBytes(br)
}

func (f FsSource) ToWasmSource() aa.WasmSource {
	return func(_ context.Context) ([]byte, error) {
		return f.ToBytes()
	}
}
