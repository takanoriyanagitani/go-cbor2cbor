package util

func ComposeErr[T, U, V any](
	f func(T) (U, error),
	g func(U) (V, error),
) func(T) (V, error) {
	return func(t T) (v V, e error) {
		u, e := f(t)
		if nil != e {
			return v, e
		}
		return g(u)
	}
}

func OkFunc[T, U any](f func(T) U) func(T) (U, error) {
	return func(t T) (U, error) {
		return f(t), nil
	}
}
