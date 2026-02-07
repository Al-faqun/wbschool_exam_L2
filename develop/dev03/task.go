package main

import (
	"fmt"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"slices"
	"strings"
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

	var sortables = []*Sortable{}
	var sortDecor func(a *Sortable, b *Sortable) int

	if options.col > 0 {
		for _, l := range lines {
			sortable := Sortable{value: "", source: l}

			for i, v := range strings.Split(l, " ") {
				if i == (options.col - 1) {
					sortable.value = v
				}
			}

			sortables = append(sortables, &sortable)
		}

		sortDecor = func(a *Sortable, b *Sortable) int {
			// todo те, у кого нет колонки, сортируются по по алфавиту и отображаются в начале списка
			// или в конце, если указан -r
			valueA := a.value
			if "" == valueA {
				valueA = a.source
			}

			valueB := b.value
			if "" == valueB {
				valueB = b.source
			}

			return defSort(valueA, valueB, col)
		}
	} else {
		for _, l := range lines {
			sortable := Sortable{value: l, source: l}
			sortables = append(sortables, &sortable)
		}

		sortDecor = func(a *Sortable, b *Sortable) int {
			return defSort(a.value, b.value, col)
		}
	}

	if !slices.IsSortedFunc(sortables, sortDecor) {
		slices.SortStableFunc(sortables, sortDecor)
	}

	return joinSortables(sortables), nil
}

func defSort(a string, b string, col *collate.Collator) int {
	// strings.Compare doesn't work, as it does not respect language rules
	return col.CompareString(a, b)
}

func joinSortables(sortables []*Sortable) string {
	lines := []string{}
	for _, sortable := range sortables {
		lines = append(lines, sortable.source)
	}
	return strings.Join(lines, SEP)
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

type Sortable struct {
	value  string
	source string
}
