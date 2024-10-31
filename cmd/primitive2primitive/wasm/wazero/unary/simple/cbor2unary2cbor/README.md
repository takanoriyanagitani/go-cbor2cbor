## Installation

Commands:
```
cmdname='primitive2primitive/wasm/wazero/unary/simple/cbor2unary2cbor'
version='v1.4.0'

go \
	install \
	-v \
	"github.com/takanoriyanagitani/go-cbor2cbor/cmd/${cmdname}@${version}"
```

## Examples

- Requirements to run this example
  - jq
  - uv(python module)
  - cbor2(python module)
  - json2arr2cbor(https://github.com/takanoriyanagitani/go-json2cbor/tree/main/cmd/json2arr2cbor)
  - asc(assembly script compiler)
  - npm(to install asc)
  - prebuilt wasm modules(use asc to build it)

Commands:
```
jq \
	-n \
	-c \
	'[
		"hw",
		42,
		42.0,
		true,
		false
	]' |
	json2arr2cbor |
	ENV_WASM_MODULES_DIR=./wasm-modules.d/out.d \
	ENV_COLUMN_INDICES=1 \
	cbor2unary2cbor |
	python3 \
		-m uv \
		tool \
		run \
		cbor2 \
		--sequence
```

Output:
```
["hw", 84, 42.0, true, false]
```
