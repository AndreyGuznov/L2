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
=== Утилита cut ===
Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные
Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const filepath = "D:/Golang/L2/dev06/text.txt"

func readData(filePath string) []string {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
	}
	sl := strings.SplitAfter(string(data), "\n")
	for _, v := range sl {
		v = strings.ReplaceAll(v, "\n", "")
		// v = strings.Replace(v, "\r", "", -1)
	}
	return sl
}

var (
	n int
	s string
	b bool
)

func main() {
	flag.IntVar(&n, "f", 0, "fields")
	flag.StringVar(&s, "d", "	", "delimiter")
	flag.Bool("s", false, "separated")
	flag.Parse()
	dataSlStr := readData(filepath)
	for i, v := range os.Args {
		if v == "-d" {
			s = os.Args[i+1]
		}
		if v == "-s" {
			b = true
		}
	}
	for _, v := range dataSlStr {

		if !b && strings.Count(v, s) > 0 {
			res := strings.Split(v, s)
			fmt.Println(res[n])
		}

	}
}
