package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
Утилита cut

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддерживать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Структура, хранящая флаги.
// Позволяет не использовать глобальные переменные.
type flags struct {
	fields    []int
	delimiter string
	separated bool
}

// Функция - конструктор.
func newFlags() *flags {
	fields := []int{}
	f := flags{
		fields: fields,
	}
	return &f
}

// Метод для парсинга -f флага.
func (f *flags) parseFFlag(strF string) error {
	//Если флаг не был передан, ничего не делаем.
	if strF == "" {
		return nil
	}
	field := []int{}
	//Проверяем было ли в поле флага перечисление.
	//И пилим на массив.
	tmpStr := strings.Split(strF, ",")
	for _, v := range tmpStr {
		switch len(v) {
		//Если перечисления не было, конвертируем стринг в инт,
		//так как это значит что в поле флага 1 число.
		case 1:
			i, err := strconv.Atoi(v)
			if err != nil {
				//log.Println("ошибка парсинга int из строки (case 1)")
				return err
			}
			field = append(field, i)
		//Если длина поля более 1:
		default:
			//Проверяем был ли передан диапазон.
			tmpStr2 := strings.Split(v, "-")
			//Пытаемся конвертировать 1 значение массива в инт.
			l, err := strconv.Atoi(tmpStr2[0])
			if err != nil {
				//log.Println("ошибка парсинга int из строки (def 0)")
				return err
			}
			//Пытаемся конвертировать 2 значение в инт.
			r, err := strconv.Atoi(tmpStr2[1])
			if err != nil {
				//log.Println("ошибка парсинга int из строки (def 1)")
				return err
			}
			//Если это диапазон, то между l и r будет разница,
			//которую надо восстановить в fields.
			for j := l; j <= r; j++ {
				field = append(field, j)
			}
		}
	}
	f.fields = field
	//log.Println("парсинг полей для вывода успешно завершен")
	return nil
}

// Метод для разделения входых строк по делителю.
func (f *flags) cut(data []string) [][]string {
	res := [][]string{}
	for _, v := range data {
		res = append(res, strings.Split(v, f.delimiter))
	}
	return res
}

// Метод для построения строки для вывода.
func (f *flags) print(rows [][]string) string {
	var b strings.Builder
	//Если стоит флаг -s, то выводим только строки с разделителем.
	if f.separated {
		for _, v := range rows {
			switch len(v) {
			case 1:
				continue
			default:
				b.WriteString("\n")
				b.WriteString(strings.Join(v, f.delimiter))
			}
		}
		return b.String()
	}
	//Если не было передано флага -f, то печатаем все,
	//собирая строки по разделителям.
	if len(f.fields) == 0 {
		for _, v := range rows {
			b.WriteString("\n")
			b.WriteString(strings.Join(v, f.delimiter))
		}
		return b.String()
	}
	//Если был передан флаг -f, то собираем только указанные колонки
	//для каждой строки.
	for _, v := range rows {
		b.WriteString("\n")
		//Не придумал как избавиться от сложности О(n^2),
		//но так как ввод подразумевается только с консоли,
		//то большие объемы текста задолбаешься печатать...
		for _, n := range f.fields {
			b.WriteString(v[n] + f.delimiter)
		}
	}
	return b.String()
}

func main() {
	//Инициализируем новый экземпляр флагов.
	f := newFlags()
	var strF string
	//Декларируем и парсим флаги.
	flag.StringVar(&strF, "f", "", "select fields")
	flag.StringVar(&f.delimiter, "d", "\t", "different delimiter (default TAB)")
	flag.BoolVar(&f.separated, "s", false, "delimited strings only")
	flag.Parse()
	//Парсим поле fields.
	err := f.parseFFlag(strF)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	//Читаем stdin до тех пор, пока не напишем "!stop"
	data := []string{}
	reader := bufio.NewScanner(os.Stdin)
	fmt.Println("Для прекращения ввода напишите !stop.")
	for reader.Scan() {
		d := reader.Text()
		if d == "!stop" {
			break
		}
		data = append(data, d)
	}
	rows := f.cut(data)
	fmt.Println(f.print(rows))
}
