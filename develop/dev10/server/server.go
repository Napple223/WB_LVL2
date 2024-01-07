package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

//Сервер для проверки работоспособности клиента.

const (
	proto = "tcp"
	addr  = ":8080"
)

func handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "read from conn err: %v\n", err)
			return
		}

		msg := strings.TrimSuffix(string(b), "\n")
		msg = strings.TrimSuffix(msg, "\r")
		_, err = conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "writing to client err: %v", err)
			return
		}
	}
}

func main() {
	listener, err := net.Listen(proto, addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "conn err: %v\n", err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "accept conn err: %v\n", err)
		}
		go handleConn(conn)
	}

}
