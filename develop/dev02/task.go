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
	var prev rune                    // руна, которая ожидает множитель
	var numBuffer strings.Builder    // буффер множителя
	var escapeMode bool

	fmt.Printf("input: %s\n", input)

	// todo: подумай, какой инпут нужен для определения действия для n-го символа в строке (сколько предыдущих символов нужно держать в памяти),
	// и напиши красивую функцию
	for _, cur := range input {
		if !escapeMode && prev == '\\' && isSpecial(cur) {
			prev = cur
			escapeMode = true
			// todo: не попасть сюда после предыдущего ифа
		} else if isNumber(cur) {
			if prev == 0 {
				return "", fmt.Errorf("String must begin with a letter") // todo: убрать дублирование ошибки?
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
		fmt.Printf("Writing rune %q\n", let)
		(*resultBuffer).WriteRune(let)
	}

	return nil
}
