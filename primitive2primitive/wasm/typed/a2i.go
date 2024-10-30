package typed

func AnyToI32(input any) (int32, error) {
	return AnyToNumber(
		input,
		func(i int32) int32 { return i },
		func(i int64) int32 { return int32(i) },
		func(i uint32) int32 { return int32(i) },
		func(i uint64) int32 { return int32(i) },
		func(i float32) int32 { return int32(i) },
		func(i float64) int32 { return int32(i) },
	)
}

func AnyToI64(input any) (int64, error) {
	return AnyToNumber(
		input,
		func(i int32) int64 { return int64(i) },
		func(i int64) int64 { return i },
		func(i uint32) int64 { return int64(i) },
		func(i uint64) int64 { return int64(i) },
		func(i float32) int64 { return int64(i) },
		func(i float64) int64 { return int64(i) },
	)
}
