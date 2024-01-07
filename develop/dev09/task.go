package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const (
	singlePage = "index.html"
)

func saveSinglePage(path string) error {
	resp, err := http.Get(path)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("getting %s: %s", path, resp.Status)
	}

	file, err := os.Create(singlePage)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func main() {
	var (
		recursion        bool
		levelOfRecursion uint64
	)

	flag.BoolVar(&recursion, "r", false, "recursively download")
	flag.Uint64Var(&levelOfRecursion, "l", 1, "recursion depth")
	flag.Parse()
	path := flag.Args()[0]
	if path == "" {
		fmt.Fprintln(os.Stderr, "no url provided")
		os.Exit(1)
	}

	switch {
	case recursion:

	default:
		fmt.Println("Starting download single page.")
		err := saveSinglePage(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse %s: %v", path, err)
			os.Exit(1)
		}
		fmt.Println("Download successfully completed.")
	}
}
