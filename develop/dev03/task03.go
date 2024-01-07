package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
Утилита sort.

Отсортировать строки в файле по аналогии с консольной утилитой sort
(man sort — смотрим описание и основные параметры):
на входе подается файл с несортированными строками,
на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:
-k — указание колонки для сортировки
(слова в строке могут выступать в качестве колонок,
по умолчанию разделитель — пробел)
-n — сортировать по числовому значению ++
-r — сортировать в обратном порядке +
-u — не выводить повторяющиеся строки

Дополнительно:
Реализовать поддержку утилитой следующих ключей:
-M — сортировать по названию месяца
-c — проверять отсортированы ли данные ++
-h — сортировать по числовому значению с учетом суффиксов ++

Базово можно сортировать только числа и буквы.
Если нет ключа -n сортировка идет по буквенным значениям.
Для буквенной и циферной сортировки можно указать все остальные ключи.
-k, -r, -u
-M только для буквенной сортировки
-h только для циферной сортировки.

Стоит ли вынести -k к базе?
*/

const (
	outputFileName = "./output_data.txt"
	perm           = 0777
)

// Структура для хранения флагов.
type flags struct {
	column    uint64
	numeric   bool
	reverse   bool
	unique    bool
	sorted    bool
	humanRead bool
}

// Структура для хранения строк.
type numeric struct {
	column int
	suffix string
	row    []string
}

// Тип, имплементирующий sort.Interface
type nums []numeric

// Методы, исплементирующие sort.Interface.
func (n nums) Len() int {
	return len(n)
}

func (n nums) Less(i, j int) bool {
	if n[i].suffix != "" && n[j].suffix != "" {
		if n[i].column <= n[j].column {
			return n[i].suffix < n[j].suffix
		}
	}
	return n[i].column < n[j].column
}

func (n nums) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// Функция для чтения строк из файла.
func readFile(inputFileName string) ([]string, error) {
	inputData, err := os.OpenFile(inputFileName, os.O_RDONLY, perm)
	defer func() {
		err := inputData.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка закрытия файла: %v\n", err)
		}
	}()
	if err != nil {
		return nil, err
	}
	d, err := io.ReadAll(inputData)
	if err != nil {
		return nil, err
	}
	res := strings.Split(string(d), "\n")
	for i, v := range res {
		res[i] = strings.TrimSuffix(v, "\r")
	}
	return res, nil
}

// Функция для разделения строки на слайс строк.
func splitData(data []string) [][]string {
	res := make([][]string, len(data))
	for i, v := range data {
		s := strings.Split(v, " ")
		res[i] = s
	}
	return res
}

// Функция для сортировки чисел.
func (f *flags) columnNumericSort(data [][]string) ([]string, error) {
	if f.column > 0 {
		for _, v := range data {
			if len(v) <= int(f.column) {
				return joinData(data),
					fmt.Errorf("номер колонки %d больше количества колонок в строке %d\n", f.column, len(v))
			}
		}
	}

	n := make([]numeric, len(data))

	//Флаг -h
	switch f.humanRead {
	case true:
		for i, v := range data {
			d, suf, err := trimSuffix(v[f.column])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			n[i] = numeric{
				column: d,
				suffix: suf,
				row:    v,
			}
		}
	default:
		for i, v := range data {
			d, err := strconv.Atoi(v[f.column])
			if err != nil {
				return nil, fmt.Errorf("ошибка парсинга цифры из строки: %v\n", err)
			}
			n[i] = numeric{
				column: d,
				suffix: "",
				row:    v,
			}
		}
	}
	sort.Sort(nums(n))

	if f.reverse {
		sort.Sort(sort.Reverse(nums(n)))
	}

	res := make([]string, len(data))
	for i, v := range n {
		res[i] = strings.Join(v.row, " ")
	}
	return res, nil
}

// Функция для обрезки буквенного суффикса для флага -h
func trimSuffix(column string) (int, string, error) {
	var sep int

	for i, v := range column {
		_, err := strconv.Atoi(string(v))
		if err == nil {
			continue
		}
		sep = i
	}
	if sep == 0 {
		return -1, "", fmt.Errorf("в колонке %s нет численной части\n", column)
	}
	d, err := strconv.Atoi(column[:sep])
	if err != nil {
		return -1, "", fmt.Errorf("ошибка парсинга численной части: %v\n", err)
	}
	suffix := column[sep:]
	return d, suffix, nil
}

// Функция для сборки слайса слайсов строк в слайс строк.
func joinData(data [][]string) []string {
	res := make([]string, len(data))

	for i, v := range data {
		s := strings.Join(v, " ")
		res[i] = s
	}
	return res
}

// Функция - точка входа в сортировку чисел.
func (f *flags) numericSort(inputFile string) ([]string, error) {
	data, err := readFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %v\n", err)
	}
	return f.columnNumericSort(splitData(data))
}

// Функция - точка входа в сортировку буквенных значений.
func (f *flags) abcSort(inputFile string) ([]string, error) {
	data, err := readFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %v\n", err)
	}
	return f.abcColumnSort(splitData(data))
}

// Функция для сортировки букв.
func (f *flags) abcColumnSort(data [][]string) ([]string, error) {
	if f.column > 0 {
		for _, v := range data {
			if len(v) <= int(f.column) {
				return joinData(data),
					fmt.Errorf("номер колонки %d больше количества колонок в строке %d\n", f.column, len(v))
			}
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i][f.column] < data[j][f.column]
	})

	if f.reverse {
		j := len(data) - 1
		for i := 0; i < j; i++ {
			data[i], data[j] = data[j], data[i]
			j--
		}
	}
	return joinData(data), nil
}

// Функция для записи отсортированных данных в файл.
func (f *flags) writeData(data []string, outputFileName string) error {
	outputFile, err := os.OpenFile(outputFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer func() {
		err := outputFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка закрытия файла для записи: %v\n", err)
		}
	}()
	if err != nil {
		return fmt.Errorf("Ошибка открытия файла для записи результатов: %v", err)
	}
	writer := bufio.NewWriter(outputFile)
	var tmp string
	for _, v := range data {
		if f.unique {
			if v == tmp {
				tmp = v
				continue
			}
		}
		_, err := writer.WriteString(v + "\n")
		if err != nil {
			return fmt.Errorf("ошибка записи данных в файл: %v", err)
		}
	}
	writer.Flush()
	return nil
}

func main() {
	//Парсим флаги.
	f := flags{}
	flag.Uint64Var(&f.column, "k", 0, "column for sorting")
	flag.BoolVar(&f.numeric, "n", false, "sort nums")
	flag.BoolVar(&f.reverse, "r", false, "reverse sorting order")
	flag.BoolVar(&f.unique, "u", false, "only unique lines")
	flag.BoolVar(&f.sorted, "c", false, "data is sorted")
	flag.BoolVar(&f.humanRead, "h", false, "sort by numeric value with suffix")
	flag.Parse()

	args := flag.Args()
	var (
		outputFile, inputFile string
	)
	//Определяем имена и пути входного и выходного файлов.
	switch len(args) {
	case 0:
		fmt.Fprintln(os.Stderr, "недостаточно аргументов.")
		os.Exit(1)
	case 1:
		inputFile = args[0]
		outputFile = outputFileName
	case 2:
		inputFile = args[0]
		outputFile = args[1]
	default:
		fmt.Fprintln(os.Stderr, "слишком много аргументов.")
		os.Exit(1)
	}

	var d []string

	switch {
	case f.sorted:
		data, err := readFile(inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		ok := sort.StringsAreSorted(data)
		fmt.Println(ok)
	case f.numeric:
		sortedData, err := f.numericSort(inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		d = sortedData
	default:
		sortedData, err := f.abcSort(inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		d = sortedData
	}

	err := f.writeData(d, outputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
