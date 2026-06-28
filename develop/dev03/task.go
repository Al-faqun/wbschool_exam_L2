package main

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

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"wbschool_exam_L2/develop/dev03/sort"
)

const SEP = "\n"

/*
* Формат запуска: ./app <flags> file_path
* Пример: ./dev03 -n extra/numsort.txt

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
/* Does not support win separators for simplicity
*/
func main() {
	// parse flags
	kPtr := flag.Int("k", 0, "Sort by column")
	nPtr := flag.Bool("n", false, "Sort as numbers")
	rPtr := flag.Bool("r", false, "Sort in reverse order")
	uPtr := flag.Bool("u", false, "Remove duplicates")

	flag.Parse()
	options := sort.SortOptions{Col: *kPtr, IsNum: *nPtr, IsRev: *rPtr, IsRemDub: *uPtr}

	// read file
	args := os.Args[1:]

	if len(args) == 0 {
		panic("No file is provided")
	}

	data, err := os.ReadFile(args[len(args)-1])
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), SEP)

	// sort
	sorted, err := sort.Sort(lines, options)

	if err != nil {
		panic(err)
	}

	// output result to stdout
	w := bufio.NewWriter(os.Stdout)
	for _, l := range sorted {
		w.WriteString(fmt.Sprintf("%s\n", l))
	}
	w.Flush()
}
