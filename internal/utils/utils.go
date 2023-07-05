package utils

import (
	"fmt"
	"strings"
)

func WrapError(text string, err error) error {
	return fmt.Errorf(text, err)
}

func Normalizetext(input string) string {
	t := strings.ToLower(strings.TrimPrefix(input, "to "))
	t = strings.TrimPrefix(t, "a ")
	t = strings.ReplaceAll(t, "`", "'")
	t = strings.Trim(t, " ")
	return t
}

func ReplaceEverySecondSymbol(input string) string {
	chars := []rune(input)
	for i := 0; i < len(chars); i++ {
		if i%2 != 0 {
			chars[i] = '*'
		}
	}
	return string(chars)
}

const (
	StartOperation = iota
	AddOperation
	TrainOperation
	LearnOperation
)
