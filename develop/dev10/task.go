package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту
(ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет,
а данные полученные из сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на
подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать
сокет и завершаться. Если сокет закрывается
со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер,
программа должна завершаться через timeout.
*/

const proto = "tcp"

func main() {
	//Парсим флаг таймаута.
	var timeout string
	flag.StringVar(&timeout, "timeout", "10s", "timeout")
	flag.Parse()
	duration, err := strconv.Atoi(strings.TrimSuffix(timeout, "s"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "wrong timeout format: got = %s, want 10s\n", timeout)
		os.Exit(1)
	}

	//Парсим хост и порт.
	host := flag.Arg(0)
	port := flag.Arg(1)

	//Задаем таймаут работы клиента.
	deadline := time.Duration(duration) * time.Second
	t := time.NewTimer(deadline)

	//Подключаемся к серверу.
	conn, err := net.Dial(
		proto,
		net.JoinHostPort(host, port),
	)
	//Если err != nil, т.е. несуществующий сервер,
	//ждем таймаут и завершаем работу.
	if err == nil {
		defer func(conn net.Conn) {
			err = conn.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "conn close err: %v\n", err)
				os.Exit(1)
			}
		}(conn)
		//Горутина для отправки сообщений на сервер и чтение от него данных.
		go func(conn net.Conn) {
			reader := bufio.NewReader(os.Stdin)
			readFromServer := bufio.NewReader(conn)
			for {
				t, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintf(os.Stderr, "reading from stdin err: %v\n", err.Error())
					os.Exit(1)
				}
				//Если ввели ^D закрываем сокет и завершаем работу клиента.
				//По идее это не должно касаться сервера.
				if t[0] == 4 {
					//Нужно явно вызвать conn.Close()
					//потому что defer не сработает.
					conn.Close()
					os.Exit(0)
				}
				//Пишем сообщение на сервер.
				_, err = conn.Write([]byte(t))
				if err != nil {
					fmt.Fprintf(os.Stderr, "writing to server err: %v\n", err)
					os.Exit(1)
				}
				//Читаем ответ.
				b, err := readFromServer.ReadBytes('\n')
				if err != nil {
					fmt.Fprintf(os.Stderr, "read response from server err: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("response from server: %s", string(b))
			}
		}(conn)
	}

	//По таймауту завершаем работу клиента.
	<-t.C
}
