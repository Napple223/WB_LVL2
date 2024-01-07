package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

/*
Задача на распаковку.
Создать Go-функцию, осуществляющую
примитивную распаковку строки,
содержащую повторяющиеся символы/руны, например:
"a4bc2d5e" => "aaaabccddddde"
"abcd" => "abcd"
"45" => "" (некорректная строка)
"" => ""

Дополнительно:
Реализовать поддержку escape-последовательностей.
Например:
qwe\4\5 => qwe45 (*)
qwe\45 => qwe44444 (*)
qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка,
функция должна возвращать ошибку. Написать unit-тесты.
*/

func unzipString(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}
	if _, err := strconv.Atoi(input); err == nil {
		return "", errors.New("invalid string")
	}
	var (
		prevSymbol rune
		escaped    bool
		b          strings.Builder
		count      int
	)

	for _, r := range input {
		if unicode.IsDigit(r) && !escaped {
			count, _ = strconv.Atoi(string(r))
			b.WriteString(strings.Repeat(string(prevSymbol), count-1))
			continue
		}
		escaped = r == '\\' && prevSymbol != '\\'
		if !escaped {
			b.WriteRune(r)
		}
		prevSymbol = r
	}
	return b.String(), nil
}
