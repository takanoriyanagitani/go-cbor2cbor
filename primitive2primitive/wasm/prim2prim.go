package prim2prim

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	p2p "github.com/takanoriyanagitani/go-cbor2cbor/primitive2primitive"
)

// Gets wasm bytes for the specified name(e.g, filename).
type WasmStore func(context.Context, string) ([]byte, error)

func WasmStoreNewFsTrusted(f fs.FS) WasmStore {
	return func(_ context.Context, moduleName string) ([]byte, error) {
		return fs.ReadFile(f, moduleName)
	}
}

func WasmStoreNewFsTrustedFromEnv(keyToTrustedWasmDirName string) WasmStore {
	var f fs.FS = os.DirFS(os.Getenv(keyToTrustedWasmDirName))
	return WasmStoreNewFsTrusted(f)
}

type FunctionInfo struct {
	// The name of the function in a wasm module.
	FunctionName string

	// The name(e.g, filename) of the wasm module.
	ModuleName string

	WasmStore
}

type IndexToFunc struct {
	// The index of the item in an array to apply the function.
	Index uint8 `json:"index"`

	// The name of the wasm module(e.g, filename).
	Module string `json:"module"`

	// The name of the function in the wasm module.
	Name string `json:"name"`

	// The type name of the function.
	FunctionType string `json:"typ"`
}

var FunctionTypeMap map[string]any = map[string]any{
	"IntToInt32":       p2p.IntToInt32(nil),
	"UintToUint32":     p2p.UintToUint32(nil),
	"FloatToFloat32":   p2p.FloatToFloat32(nil),
	"IntToInt64":       p2p.IntToInt64(nil),
	"UintToUint64":     p2p.UintToUint64(nil),
	"FloatToFloat64":   p2p.FloatToFloat64(nil),
	"LongToDouble":     p2p.LongToDouble(nil),
	"TimeFromUnixtime": p2p.TimeFromUnixtime(nil),
}

// Converts the name of the type to the type using [FunctionTypeMap].
func (i IndexToFunc) ToFunctionType() any {
	var tname string = i.FunctionType
	typ, found := FunctionTypeMap[tname]
	switch found {
	case true:
		return typ
	default:
		return fmt.Errorf("invalid type name: %s", tname)
	}
}

type IndexToFunction interface {
	ToFunctionArray(context.Context) ([]p2p.PrimitiveToPrimitiveA, error)
	Close(context.Context) error
}

type IndexToFuncStore func(context.Context) ([]IndexToFunc, error)

type IndexToFuncStoreRaw func(context.Context) ([]byte, error)

func IndexToFuncStoreRawFsNew(
	f fs.FS,
	trustedFilename string,
) IndexToFuncStoreRaw {
	return func(_ context.Context) ([]byte, error) {
		return fs.ReadFile(f, trustedFilename)
	}
}

func IndexToFuncStoreRawFsNewFromEnv(
	keyToTrustedDirName string,
	keyToTrustedCfgName string,
) IndexToFuncStoreRaw {
	var f fs.FS = os.DirFS(os.Getenv(keyToTrustedDirName))
	return IndexToFuncStoreRawFsNew(
		f,
		os.Getenv(keyToTrustedCfgName),
	)
}

type IndexToFuncParser func([]byte) ([]IndexToFunc, error)

func IndexToFuncParserJson(raw []byte) (parsed []IndexToFunc, e error) {
	e = json.Unmarshal(raw, &parsed)
	return parsed, e
}

func (r IndexToFuncStoreRaw) ToStore(
	parser IndexToFuncParser,
) IndexToFuncStore {
	return func(ctx context.Context) ([]IndexToFunc, error) {
		raw, e := r(ctx)
		if nil != e {
			return nil, e
		}
		return parser(raw)
	}
}

func (r IndexToFuncStoreRaw) ToJsonStore() IndexToFuncStore {
	return r.ToStore(IndexToFuncParserJson)
}
