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

func Values[T comparable, S any](m map[T]S) (s []S) {
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

func StripContent(content string, contentMaxSize int) string {
	contentRune := []rune(content)
	contentRuneLength := len(contentRune)
	if contentRuneLength <= contentMaxSize {
		return content
	}
	return string(contentRune[:contentMaxSize])
}
