package main

import (
	"slices"
	"strings"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const SEP = "\n"

/* Does not support windows separators for simplicity */
func main() {

}

func sort(input string, options SortOptions) (string, error) {
	// without options
	lines := strings.Split(input, SEP)

	col := createCol(options)

	sortDec := func (a string, b string) int {
		return defSort(a, b, col)
	}

	if !slices.IsSortedFunc(lines, sortDec) {
		slices.SortStableFunc(lines, sortDec)
	}

	return strings.Join(lines, SEP), nil
}

func defSort(a string, b string, col *collate.Collator) int {
	// strings.Compare doesn't work, as it does not respect language rules
	return strings.Compare(a, b)
}

func createCol(options SortOptions) *collate.Collator {
	// todo: for numeric collator see https://pkg.go.dev/golang.org/x/text/collate#example-New
	return collate.New(language.English)
}

type SortOptions struct {
	col      int
	isNum    bool
	isRev    bool
	isRemDub bool
}
