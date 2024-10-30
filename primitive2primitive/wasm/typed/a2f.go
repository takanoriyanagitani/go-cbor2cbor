package typed

func AnyToF64(input any) (float64, error) {
	return AnyToNumber(
		input,
		func(i int32) float64 { return float64(i) },
		func(i int64) float64 { return float64(i) },
		func(i uint32) float64 { return float64(i) },
		func(i uint64) float64 { return float64(i) },
		func(i float32) float64 { return float64(i) },
		func(i float64) float64 { return i },
	)
}

func AnyToF32(input any) (float32, error) {
	return AnyToNumber(
		input,
		func(i int32) float32 { return float32(i) },
		func(i int64) float32 { return float32(i) },
		func(i uint32) float32 { return float32(i) },
		func(i uint64) float32 { return float32(i) },
		func(i float32) float32 { return i },
		func(i float64) float32 { return float32(i) },
	)
}
