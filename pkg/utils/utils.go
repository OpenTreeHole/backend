package utils

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

func StripContent(content string, contentMaxSize int) string {
	contentRune := []rune(content)
	contentRuneLength := len(contentRune)
	if contentRuneLength <= contentMaxSize {
		return content
	}
	return string(contentRune[:contentMaxSize])
}
