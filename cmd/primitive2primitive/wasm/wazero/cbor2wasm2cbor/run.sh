#!/bin/sh

export ENV_TRUSTED_WASM_DIR=./sample.d/modules.d

export ENV_TRUSTED_CFG_DIR_NAME=./sample.d
export ENV_TRUSTED_CFG_FILENAME=empty.json
export ENV_TRUSTED_CFG_FILENAME=fcfg.json

./simple
