package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
=== Утилита grep ===
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

var (
	n int
	s string
)

const (
	filePath = "D:/Golang/L2/Task5/text.txt"
)

func outAfter(n int, sl []string) {
	for i := 0; i < len(sl); i++ {
		sl[i] = strings.ReplaceAll(sl[i], "\r", "")
		if sl[i] == os.Args[3] {
			for j := i + 1; j <= i+n; j++ {
				fmt.Println(sl[j])
			}
		}
	}
}

func outBefore(n int, sl []string) {
	for i := 0; i < len(sl); i++ {
		sl[i] = strings.ReplaceAll(sl[i], "\r", "")
		if sl[i] == os.Args[3] {
			for j := i - n - 1; j < n; j++ {
				fmt.Println(sl[j])
			}
		}
	}
}

func outInvert(n int, sl []string) {
	copy(sl[n:], sl[n+1:])
	for _, val := range sl {
		fmt.Println(val)
	}
}

func outFixed(s string, sl []string) {
	for i := 0; i < len(sl); i++ {
		sl[i] = strings.ReplaceAll(sl[i], "\r", "")
		if sl[i] == s {
			fmt.Println(sl[i])
		}
	}
}

func outLineNum(s string, sl []string) {
	for i := 0; i < len(sl); i++ {
		sl[i] = strings.ReplaceAll(sl[i], "\r", "")
		if sl[i] == s {
			fmt.Println(i)
		}
	}
}

func main() {
	flag.IntVar(&n, "A", 0, "after")
	flag.IntVar(&n, "B", 0, "before")
	flag.IntVar(&n, "C", 0, "context")
	flag.String("c", "0", "count")
	flag.IntVar(&n, "i", 0, "ignore-case")
	flag.IntVar(&n, "v", 0, "invert")
	flag.StringVar(&s, "F", " ", "fixed")
	flag.StringVar(&s, "A", "0", "line num")

	flag.Parse()
	// fmt.Println(os.Args[1])
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
	}
	dataSlStr := strings.Split(string(data), "\n")

	switch os.Args[1] {
	case "-A":
		outAfter(n, dataSlStr)
	case "-B":
		outBefore(n, dataSlStr)
	case "-C":
		outBefore(n, dataSlStr)
		outAfter(n, dataSlStr)
	case "-c":
		fmt.Println(len(dataSlStr))
	case "i":
		fmt.Println(dataSlStr[n])
	case "v":
		outInvert(n, dataSlStr)
	case "F":
		outFixed(s, dataSlStr)
	case "n":
		outLineNum(s, dataSlStr)
	default:
		fmt.Println("You are thinking that I am a Wizard???")
	}
}
