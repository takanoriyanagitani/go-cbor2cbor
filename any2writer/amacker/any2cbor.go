package any2cbor

import (
	"io"
	"os"

	ac "github.com/fxamacker/cbor/v2"

	a2w "github.com/takanoriyanagitani/go-cbor2cbor/any2writer"
)

func TimeModeFromString(timeMode string) ac.TimeMode {
	switch timeMode {
	case "TimeUnix":
		return ac.TimeUnix
	case "TimeUnixMicro":
		return ac.TimeUnixMicro
	case "TimeUnixDynamic":
		return ac.TimeUnixDynamic
	case "TimeRFC3339":
		return ac.TimeRFC3339
	case "TimeRFC3339Nano":
		return ac.TimeRFC3339Nano
	default:
		return ac.TimeUnix
	}
}

var TimeModeFromEnv ac.TimeMode = TimeModeFromString(
	os.Getenv("ENV_CBOR_TIME_MODE"),
)

type AnyArrayToCbor struct {
	*ac.Encoder
}

func (a AnyArrayToCbor) Encode(arr []any) error {
	return a.Encoder.Encode(arr)
}

func AnyArrayToCborNew(wtr io.Writer) func([]any) error {
	var eopt ac.EncOptions = ac.CanonicalEncOptions()
	eopt.Time = TimeModeFromEnv
	em, err := eopt.EncMode()
	var a2c AnyArrayToCbor
	if nil == err {
		a2c.Encoder = em.NewEncoder(wtr)
	}
	return func(arr []any) error {
		if nil != err {
			return err
		}
		return a2c.Encode(arr)
	}
}

var AnyToWriterNew a2w.AnyToWriterNew = AnyArrayToCborNew
