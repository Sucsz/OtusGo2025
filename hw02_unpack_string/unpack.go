package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var builder strings.Builder

	for i := 0; i < len(str); {
		char, size := utf8.DecodeRuneInString(str[i:])
		// strconv.Atoi(...) делает код не очень читаемым из-за принципа возвращения ошибки
		// Правда unicode.IsDigit(...) кушает и не арабские цифры, как цифры
		if unicode.IsDigit(char) {
			// если символ цифра, то она или первая в строке или предыдущий символ тоже цифра
			return "", ErrInvalidString
		}
		if i+size < len(str) {
			// следующий символ существует
			nextChar, nextSize := utf8.DecodeRuneInString(str[i+size:])
			if unicode.IsDigit(nextChar) {
				builder.WriteString(strings.Repeat(string(char), int(nextChar-'0')))
				i += size + nextSize
				continue
			}
		}
		builder.WriteRune(char)
		i += size
	}
	return builder.String(), nil
}
