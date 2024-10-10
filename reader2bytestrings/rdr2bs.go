package rdr2bs

import (
	"io"
)

type RdrToByteStrs func(io.Reader) func(*[][]byte) error
