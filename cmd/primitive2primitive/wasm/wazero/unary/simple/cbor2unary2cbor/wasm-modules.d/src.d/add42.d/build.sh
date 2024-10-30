#!/bin/sh

inpt=./add42.ts
out=./add42.wasm

asc \
	--outFile "${out}" \
	--target release \
	-Osize \
	--optimizeLevel 3 \
	--shrinkLevel s \
	--converge \
	--enable simd \
	--enable relaxed-simd \
	"${inpt}"
