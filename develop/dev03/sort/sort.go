package sort

import (
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const SEP = "\n"

func Sort(input string, options SortOptions) (string, error) {
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

	if options.IsRemDub == true {
		sortables = slices.CompactFunc(sortables, func(a *Sortable, b *Sortable) bool {
			return strings.EqualFold(a.Source, b.Source)
		})
	}

	return joinSortables(sortables), nil
}

func prepare(lines []string, options SortOptions) ([]*Sortable, error) {
	var sortables = []*Sortable{}

	for _, l := range lines {
		sortable := Sortable{Value: l, Source: l}

		sortables = append(sortables, &sortable)
	}

	if options.Col > 0 {
		for _, sortable := range sortables {
			sortable.Value = ""

			for i, v := range strings.Split(sortable.Source, " ") {
				if i == (options.Col - 1) {
					sortable.Value = v
					break
				}
			}
		}
	}

	if options.IsNum == true {
		re, err := regexp.Compile(`^\s*(-?\d[\d.]*)`)
		if err != nil {
			return nil, err
		}

		for _, sortable := range sortables {
			matches := re.FindStringSubmatch(sortable.Value)
			if matches != nil {
				sortable.Value = matches[1]
			} else {
				sortable.Value = ""
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

	if options.IsNum == true {
		innerFunc = numSort
	} else {
		innerFunc = func(a *Sortable, b *Sortable, col *collate.Collator) int {
			return defSort(a.Value, b.Value, col)
		}
	}

	if options.IsRev {
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

	emptyA := len(a.Value) == 0
	emptyB := len(b.Value) == 0

	if !emptyA {
		aNum, err = strconv.ParseFloat(a.Value, 64)
		if err != nil {
			panic(err)
		}
	}

	if !emptyB {
		bNum, err = strconv.ParseFloat(b.Value, 64)
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
		return defSort(a.Source, b.Source, col)
	}

	// both are numbers
	if aNum < bNum {
		return -1
	} else if aNum > bNum {
		return 1
	} else {
		// if both numbers are equal, their full strings are sorted alphabetically
		return defSort(a.Source, b.Source, col)
	}
}

func joinSortables(sortables []*Sortable) string {
	lines := []string{}
	for _, sortable := range sortables {
		lines = append(lines, sortable.Source)
	}
	return strings.Join(lines, SEP)
}

func createCol(options SortOptions) *collate.Collator {
	// todo: for numeric collator see https://pkg.go.dev/golang.org/x/text/collate#example-New
	return collate.New(language.English)
}

type SortOptions struct {
	Col      int
	IsNum    bool
	IsRev    bool
	IsRemDub bool
}

type Sortable struct {
	Value  string
	Source string
}
