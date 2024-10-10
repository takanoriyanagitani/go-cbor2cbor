package any2writer

import (
	"io"
)

type AnyToWriterNew func(io.Writer) func([]any) error
