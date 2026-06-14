package main

import (
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"regexp"
	"slices"
	"strconv"
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

/*
* Если указаны четыре ключа (-k, -n, -r, -u)
* k: pre + algo
* n: pre + algo
* r: algo
* u: post
*
* Порядок выполнения флагов:
* 1. k - получаем значения из колонки
* 2. n - сортируем выбранные значения как числа
* 3. r - какой бы алгоритм мы не выбрали, применяем к нему реверс
* 4. u - убираем дубликаты после сортировки (применяем фильтр к сортируемому значению)
*
/* Does not support windows separators for simplicity
*/
func main() {
	// todo parse flags
	// read file
	// sort
	// output result to stdout
}

func sort(input string, options SortOptions) (string, error) {
	lines := strings.Split(input, SEP)

	col := createCol(options)

	sortables, err := prepare(lines, options)
	if err != nil {
		return "", err
	}

	sortFunc, err := getSortFunc(options, col)
	if err != nil {
		return "", err
	}

	if !slices.IsSortedFunc(sortables, sortFunc) {
		slices.SortStableFunc(sortables, sortFunc)
	}

	if options.isRemDub == true {
		sortables = slices.CompactFunc(sortables, func(a *Sortable, b *Sortable) bool {
			return strings.EqualFold(a.source, b.source)
		})
	}

	return joinSortables(sortables), nil
}

func prepare(lines []string, options SortOptions) ([]*Sortable, error) {
	var sortables = []*Sortable{}

	for _, l := range lines {
		sortable := Sortable{value: l, source: l}

		sortables = append(sortables, &sortable)
	}

	if options.col > 0 {
		for _, sortable := range sortables {
			sortable.value = ""

			for i, v := range strings.Split(sortable.source, " ") {
				if i == (options.col - 1) {
					sortable.value = v
					break
				}
			}
		}
	}

	if options.isNum == true {
		re, err := regexp.Compile(`^\s*(-?\d[\d.]*)`)
		if err != nil {
			return nil, err
		}

		for _, sortable := range sortables {
			matches := re.FindStringSubmatch(sortable.value)
			if matches != nil {
				sortable.value = matches[1]
			} else {
				sortable.value = ""
			}
		}
	}

	return sortables, nil
}

// For ascending sort return
// -1 when a < b,
// 1 when a > b
// and zero when a == b or a and b are incomparable in the sense of a strict weak ordering.
func getSortFunc(options SortOptions, col *collate.Collator) (func(*Sortable, *Sortable) int, error) {
	var innerFunc func(a *Sortable, b *Sortable, col *collate.Collator) int
	var sortFunc func(a *Sortable, b *Sortable) int

	if options.isNum == true {
		innerFunc = numSort
	} else {
		innerFunc = func(a *Sortable, b *Sortable, col *collate.Collator) int {
			return defSort(a.value, b.value, col)
		}
	}

	if options.isRev {
		sortFunc = func(a *Sortable, b *Sortable) int {
			return innerFunc(b, a, col)
		}
	} else {
		sortFunc = func(a *Sortable, b *Sortable) int {
			return innerFunc(a, b, col)
		}
	}

	return sortFunc, nil
}

// returns 0 if a == b, -1 if a < b, and +1 if a > b
func defSort(a string, b string, col *collate.Collator) int {
	// strings.Compare doesn't work, as it does not respect language rules
	return col.CompareString(a, b)
}

// Sorts strings beginning with numbers numerically (i.e. 10 < 100)
// strings not beginning with numbers and text after numbers are sorted alphabetically,
// and are treated as 0 (placed after 0 and before the next number).
// Whitespace at the start of the string is ignored.
// Decimals are not supported.
func numSort(a *Sortable, b *Sortable, col *collate.Collator) int {
	var aNum, bNum float64
	var err error

	emptyA := len(a.value) == 0
	emptyB := len(b.value) == 0

	if !emptyA {
		aNum, err = strconv.ParseFloat(a.value, 64)
		if err != nil {
			panic(err)
		}
	}

	if !emptyB {
		bNum, err = strconv.ParseFloat(b.value, 64)
		if err != nil {
			panic(err)
		}
	}

	// strings, compared to numbers, are equal to 0.
	// strings, compared to 0, are greater than 0.
	if !emptyA && emptyB {
		// a is number, b is not a number

		if aNum <= 0 {
			return -1
		} else {
			return 1
		}
	}
	if !emptyB && emptyA {
		// b is number, a is not a number
		if bNum <= 0 {
			return 1
		} else {
			return -1
		}
	}

	// string, compared to another string, is sorted alphabetically
	if emptyA && emptyB {
		return defSort(a.source, b.source, col)
	}

	// both are numbers
	if aNum < bNum {
		return -1
	} else if aNum > bNum {
		return 1
	} else {
		// if both numbers are equal, their full strings are sorted alphabetically
		return defSort(a.source, b.source, col)
	}
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
