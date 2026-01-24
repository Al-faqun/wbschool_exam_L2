package main

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "1 argument is required")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "Only 1 argument is allowed")
		os.Exit(1)
	}

	input := os.Args[1]
	result, err := unpack(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during unpacking the input: %s", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", result)
}

/** Ограничение: функция не работает корректно с текстом, состоящим из более чем одной руны (например, à)*/
func unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	var resultBuffer strings.Builder // обработанная строка
	var prev rune                    // руна, которая ожидает множитель
	var numBuffer strings.Builder    // буффер множителя
	var escapeMode bool

	for _, cur := range input {
		if !escapeMode && prev == '\\' && isSpecial(cur) {
			prev = cur
			escapeMode = true
		} else if isNumber(cur) {
			if prev == 0 {
				return "", fmt.Errorf("String must begin with a letter")
			}
			numBuffer.WriteRune(cur)
		} else if prev == 0 {
			// записываем текущую руну в буффер
			prev = cur
		} else {
			// в буффере уже есть руна, и текущая руна не число - значит можем спокойно записать буффер в строку
			err := addRune(prev, &numBuffer, &resultBuffer)

			if err != nil {
				return "", err
			}

			escapeMode = false
			// запоминаем следующую руну
			prev = cur
		}
	}

	if prev != 0 {
		err := addRune(prev, &numBuffer, &resultBuffer)
		if err != nil {
			return "", err
		}
	}

	return resultBuffer.String(), nil
}

func isNumber(r rune) bool {
	return r >= 0x30 && r <= 0x39
}

func isSpecial(r rune) bool {
	return isNumber(r) || r == '\\'
}

func addRune(let rune, numBuffer, resultBuffer *strings.Builder) error {
	var num int // множитель
	var err error

	if (*numBuffer).Len() == 0 {
		num = 1
	} else {
		num, err = strconv.Atoi((*numBuffer).String())

		if err != nil {
			return err
		}

		(*numBuffer).Reset()
	}

	if (*resultBuffer).Len() == 0 && isNumber(let) {
		return fmt.Errorf("String must begin with a letter")
	}

	for i := 0; i < num; i++ {
		(*resultBuffer).WriteRune(let)
	}

	return nil
}
