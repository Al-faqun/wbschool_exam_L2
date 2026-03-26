package main

import (
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
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
* 4. u - убираем дубликаты после сортировки
*
/* Does not support windows separators for simplicity
*/
func main() {

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
		// Sorts strings beginning with numbers numerically (i.e. 10 < 100)
		// strings not beginning with numbers and text after numbers is sorted alphabetically
		// and are treated as 0 (placed after 0 and before the next number).
		// Decimals are not supported.
		re, err := regexp.Compile(`^\s*(-?\d+)`)
		if err != nil {
			return nil, err
		}

		for _, sortable := range sortables {
			matches := re.FindStringSubmatch(sortable.value)
			if matches != nil {
				sortable.value = matches[1]
			} else {
				sortable.value = sortable.source
			}
		}
	}

	return sortables, nil
}

func getSortFunc(options SortOptions, col *collate.Collator) (func(*Sortable, *Sortable) int, error) {
	var innerFunc func(a string, b string, col *collate.Collator) int
	var sortFunc func(a *Sortable, b *Sortable) int

	if options.isNum == true {
		innerFunc = numSort
	} else {
		innerFunc = defSort
	}

	if options.isRev {
		sortFunc = func(a *Sortable, b *Sortable) int {
			return innerFunc(b.value, a.value, col)
		}
	} else {
		sortFunc = func(a *Sortable, b *Sortable) int {
			return innerFunc(a.value, b.value, col)
		}
	}

	return sortFunc, nil
}

// todo: reverse flag
// returns 0 if a == b, -1 if a < b, and +1 if a > b
func defSort(a string, b string, col *collate.Collator) int {
	// strings.Compare doesn't work, as it does not respect language rules
	return col.CompareString(a, b)
}

func numSort(a string, b string, col *collate.Collator) int {
	// string a or b is empty?
	if len(a) == 0 {
		a = "0"
	}
	if len(b) == 0 {
		b = "0"
	}

	firstRuneA, _ := utf8.DecodeRuneInString(a)
	firstRuneB, _ := utf8.DecodeRuneInString(b)

	// string's first char is not a number?
	// if a is not a number and b is a number, b > 0
	if !unicode.IsNumber(firstRuneA) && unicode.IsNumber(firstRuneB) {
		return -1
	}
	if !unicode.IsNumber(firstRuneB) && unicode.IsNumber(firstRuneA) {
		return 1
	}
	if !unicode.IsNumber(firstRuneA) && !unicode.IsNumber(firstRuneB) {
		return defSort(a, b, col)
	}

	// both are numbers
	aNum, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	bNum, err := strconv.Atoi(b)
	if err != nil {
		panic(err)
	}

	if aNum < bNum {
		return -1
	} else if aNum > bNum {
		return +1
	} else {
		return 0
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
