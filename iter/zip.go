package itools

import (
	"iter"
	"slices"
)

type Pair[L, R any] struct {
	Left  L
	Right R
}

func Zip[L, R any](left iter.Seq[L], right iter.Seq[R]) iter.Seq[Pair[L, R]] {
	return func(yield func(Pair[L, R]) bool) {
		nl, sl := iter.Pull(left)
		nr, sr := iter.Pull(right)
		defer sl()
		defer sr()

		for {
			l, okl := nl()
			r, okr := nr()

			var ok bool = okl && okr
			var ng bool = !ok
			if ng {
				return
			}

			p := Pair[L, R]{Left: l, Right: r}

			if !yield(p) {
				sr()
				sl()
				return
			}
		}
	}
}

type Iterator[T any] iter.Seq[T]

func (i Iterator[T]) Collect() []T {
	return slices.Collect(iter.Seq[T](i))
}

func (i Iterator[T]) Zip(other Iterator[T]) iter.Seq[Pair[T, T]] {
	return Zip(
		iter.Seq[T](i),
		iter.Seq[T](other),
	)
}

func (i Iterator[T]) First(alt T) T {
	for t := range i {
		return t
	}
	return alt
}

func (i Iterator[T]) FirstOrDefault() T {
	var t T
	return i.First(t)
}

func Map[T, U any](i iter.Seq[T], mapper func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for t := range i {
			var mapd U = mapper(t)
			if !yield(mapd) {
				return
			}
		}
	}
}

func CollectErr[T any](i iter.Seq[Pair[error, T]]) ([]T, error) {
	var ret []T
	for pair := range i {
		if nil != pair.Left {
			return nil, pair.Left
		}
		ret = append(ret, pair.Right)
	}
	return ret, nil
}
