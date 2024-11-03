#!/bin/sh

export ENV_WASM_BYTES_MAX=16777216

export ENV_WASM_MODULES_DIR=./modules.d/copy/rs-copy
export ENV_WASM_FILENAME=rs_copy.wasm

export ENV_WASM_MODULES_DIR=./modules.d/copy/rs-serialized-parsed
export ENV_WASM_FILENAME=rs_serialized_parsed.wasm

jq \
	-c \
	-n \
	'[
		"hw",
		42,
		42.195,
		true,
		false
	]' |
	json2arr2cbor |
	./cbor2wasm2cbor |
	python3 \
		-m uv \
		tool \
		run \
		cbor2 \
		--sequence
