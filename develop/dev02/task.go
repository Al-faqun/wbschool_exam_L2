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
	"strconv"
	"strings"
)

func main() {

}

/** Ограничение: функция не работает корректно с текстом, состоящим из более чем одной руны (например, à)*/
func unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	var resultBuffer strings.Builder // обработанная строка
	var let rune                     // руна, которая ожидает множитель
	var numBuffer strings.Builder    // буффер множителя

	for _, r := range input {
		if isNumber(r) {
			if let == 0 {
				return "", fmt.Errorf("String must begin with a letter")
			}
			// записываем множитель
			numBuffer.WriteRune(r)
		} else if let == 0 {
			// записываем текущую руну в буффер
			let = r
		} else {
			// в буффере уже есть руна, и текущая руна не число - значит можем спокойно записать буффер в строку
			err := addRune(let, &numBuffer, &resultBuffer)

			if err != nil {
				return "", nil
			}

			// запоминаем следующую руну
			let = r
		}
	}

	if let != 0 {
		err := addRune(let, &numBuffer, &resultBuffer)
		if err != nil {
			return "", nil
		}
	}

	return resultBuffer.String(), nil
}

func isNumber(r rune) bool {
	return r >= 0x30 && r <= 0x39
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

	for i := 0; i < num; i++ {
		(*resultBuffer).WriteRune(let)
	}

	return nil
}
