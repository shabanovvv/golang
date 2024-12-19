package unpack

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(strInput string) (string, error) {
	runeInput := []rune(strInput)
	var runeOutput []rune
	var ecranLetter string

	for i, r := range runeInput {
		// Error! проверим первую руну на цифру
		if i == 0 && unicode.IsDigit(r) {
			return "", ErrInvalidString
		} else if i == 0 {
			continue
		}
		// Error! число
		if unicode.IsDigit(r) && unicode.IsDigit(runeInput[i-1]) {
			return "", ErrInvalidString
		}
		// Success! Repeat! текущая руна цифра, а предыдущая руна символ или экран
		if unicode.IsDigit(r) {
			repeatCount, _ := strconv.Atoi(string(r))
			var repeatRunes string
			if len(ecranLetter) > 0 {
				repeatRunes = strings.Repeat(ecranLetter, repeatCount)
			} else if unicode.IsDigit(r) && !unicode.IsDigit(runeInput[i-1]) {
				repeatRunes = strings.Repeat(string(runeInput[i-1]), repeatCount)
			}
			runeOutput = append(runeOutput, []rune(repeatRunes)...)
		}
		// Success! текущая руна символ или экранирование и предыдущая руна символ
		if (unicode.IsLetter(r) || r == '\\') && unicode.IsLetter(runeInput[i-1]) {
			runeOutput = append(runeOutput, runeInput[i-1])
		}
		// Success! последняя руна символ в слайсе
		if unicode.IsLetter(r) && (len(runeInput)-1) == i {
			runeOutput = append(runeOutput, r)
		}
		// определим экран
		if runeInput[i-1] == '\\' && unicode.IsLetter(r) {
			ecranLetter = string('\\') + string(r)
		} else {
			ecranLetter = ""
		}
	}

	return string(runeOutput), nil
}
