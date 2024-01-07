package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

/*
Утилита grep

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const (
	perm = 0777
)

// Структура grep.
type grep struct {
	before       uint64
	after        uint64
	context      uint64
	count        bool
	ignore       bool
	invert       bool
	fixed        bool
	lineNum      bool
	rows         [][]byte
	grepRowsNums []int
	rgx          *regexp.Regexp
}

// Функция - конструктор для инициализации массивов в grep.
func newGrep() *grep {
	rows := [][]byte{}
	rowsNums := []int{}
	g := grep{
		rows:         rows,
		grepRowsNums: rowsNums,
	}
	return &g
}

// Функция для парсинга аргументов. В зависимости от переданных
// флагов формирует регулярное выражение для поиска в тексте.
func (g *grep) parseArgs(args []string) (string, error) {
	l := len(args)
	fixBefore := `\Q`
	fixAfter := `\E`
	ignBefore := `(?i)(`
	ignAfter := `)`
	switch l {
	case 0:
		return "", errors.New("can't find pattern")
	case 1:
		regExp := args[0]
		if g.fixed {
			regExp = fixBefore + regExp + fixAfter
		}
		if g.ignore {
			regExp = ignBefore + regExp + ignAfter
		}
		rgx, err := regexp.Compile(regExp)
		if err != nil {
			return "", err
		}
		g.rgx = rgx
		return "", nil
	default:
		pipe := `|`
		regExp := args[0]
		defBefore := `(`
		defAfter := `)`
		if g.fixed {
			defBefore += fixBefore
			defAfter = fixAfter + defAfter
		}
		regExp = defBefore + regExp + defAfter
		for i := 1; i < l-1; i++ {
			regExp += pipe + defBefore + args[i] + defAfter
		}
		if g.ignore {
			regExp = ignBefore + regExp + ignAfter
		}
		rgx, err := regexp.Compile(regExp)
		if err != nil {
			return "", err
		}
		g.rgx = rgx
		return args[l-1], nil
	}
}

// Функция для добавления в массив номера строк,
// где было найдено совпадение.
func (g *grep) filterRows() {
	rowsNums := []int{}
	for idx, row := range g.rows {
		if g.rgx.Match(row) != g.invert {
			rowsNums = append(rowsNums, idx)
		}
	}
	g.grepRowsNums = rowsNums
}

// Функция для вывода результатов поиска.
func (g *grep) printRes() {
	//Если стоит флаг -c, то выведет количество строк.
	if g.count {
		fmt.Println(len(g.grepRowsNums))
		return
	}
	var (
		left  int
		right int
	)
	l := len(g.rows) - 1
	//Так как существует вероятность использования
	//флагов -A, -B, -C вместе, то необходимо
	//использовать самые широкие границы для
	//печати строк.
	if g.before < g.context {
		g.before = g.context
	}
	if g.after < g.context {
		g.after = g.context
	}

	for _, rowNum := range g.grepRowsNums {
		left = rowNum - int(g.before)
		if left < 0 {
			left = 0
		}
		right = rowNum + int(g.after)
		if right > l {
			right = l
		}
		if g.lineNum {
			for i := left; i <= right; i++ {
				if i == rowNum {
					fmt.Printf("target! %d: %s\n", rowNum, string(g.rows[i]))
					continue
				}
				fmt.Printf("%d: %s\n", rowNum, string(g.rows[i]))
			}
		} else {
			for i := left; i <= right; i++ {
				fmt.Println(string(g.rows[i]))
			}
		}
	}
}

// Структура сканера строк.
type scanner struct {
	sc   *bufio.Scanner
	file *os.File
}

// Функция - конструктор для выбора источника строк для поиска.
func newScanner(fileName string) (*scanner, error) {
	if fileName == "" {
		sc := bufio.NewScanner(os.Stdin)
		return &scanner{sc: sc}, nil
	}
	file, err := os.OpenFile(fileName, os.O_RDONLY, perm)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(file)

	scan := scanner{
		sc:   sc,
		file: file,
	}
	return &scan, nil
}

// Функция, считывающая строки в слайс слайсов байт.
func (s *scanner) readRows() [][]byte {
	rows := [][]byte{}
	for s.sc.Scan() {
		rows = append(rows, s.sc.Bytes())
	}
	return rows
}

func main() {
	//Инициализируем новый экземмпляр структуры grep.
	grep := newGrep()
	//Декларируем флаги.
	flag.Uint64Var(&grep.after, "A", 0, "print +N lines after match")
	flag.Uint64Var(&grep.before, "B", 0, "print +N lines before match")
	flag.Uint64Var(&grep.context, "C", 0, "print ±N lines around the match")
	flag.BoolVar(&grep.count, "c", false, "quantity of lines")
	flag.BoolVar(&grep.ignore, "i", false, "ignore case")
	flag.BoolVar(&grep.invert, "v", false, "exclude")
	flag.BoolVar(&grep.fixed, "F", false, "exact string match, not a pattern")
	flag.BoolVar(&grep.lineNum, "n", false, "print a number of line")
	flag.Parse()
	fileName, err := grep.parseArgs(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	sc, err := newScanner(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	grep.rows = sc.readRows()
	sc.file.Close()
	grep.filterRows()
	grep.printRes()
}
