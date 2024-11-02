package arr2cbor

import (
	"bytes"
	"context"
)

type ArrayToCbor func([]any) error

type ArrayToCborToBuffer func(context.Context, []any, *bytes.Buffer) error
