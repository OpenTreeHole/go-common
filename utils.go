package common

import "golang.org/x/exp/constraints"

type Map = map[string]any

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Keys[T comparable, S any](m map[T]S) (s []T) {
	for k := range m {
		s = append(s, k)
	}
	return s
}
