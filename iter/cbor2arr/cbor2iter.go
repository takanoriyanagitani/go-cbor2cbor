package cbor2arr

import (
	"iter"
)

type CborToArrays func() iter.Seq[[]any]
