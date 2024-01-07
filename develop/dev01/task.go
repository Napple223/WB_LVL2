package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

/*

Создать программу печатающую точное время с использованием NTP -библиотеки.
Инициализировать как go module. Использовать библиотеку github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Требования:
Программа должна быть оформлена как go module
Программа должна корректно обрабатывать ошибки библиотеки:
выводить их в STDERR и возвращать ненулевой код выхода в OS
*/

const (
	address = "0.beevik-ntp.pool.ntp.org"
)

func getNTPTime(address string) (time.Time, error) {
	time, err := ntp.Time(address)
	if err != nil {
		return time, err
	}
	return time, nil
}

func main() {
	t, err := getNTPTime(address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Println(t)
}
