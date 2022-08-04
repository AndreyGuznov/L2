package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

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
var n int

func main() {
	args := os.Args
	if len(args) == 1 {
		log.Println("Insert correct command")
	}
	argsT := strings.Join(args[1:], "")
	comand := argsT[:2]
	if len(args) > 2 {
		num := args[2:]
		n, _ = strconv.Atoi(num[0])
	}

	data, err := ioutil.ReadFile("D:/Golang/L2/Task3/text.txt")
	if err != nil {
		fmt.Println(err)
	}
	res := ""
	dataSlStr := strings.Split(string(data), "\n")
	switch comand {
	case "-k":
		for _, val := range dataSlStr {
			newVal := strings.SplitAfter(val, " ")
			newVal[0], newVal[n-1] = newVal[n-1], newVal[0]
			sort.Sort(sort.StringSlice(newVal))
			for i := 0; i < len(newVal); i++ {
				res += newVal[i] + "\n"
			}
		}
		ioutil.WriteFile("D:/Golang/L2/Task3/res.txt", []byte(res), 0777)
	case "-n":
		sort.Sort(sort.StringSlice(dataSlStr))
		for _, val := range dataSlStr {
			res += val
		}
		ioutil.WriteFile("D:/Golang/L2/Task3/res.txt", []byte(res), 0777)
	case "-r":
		for i := len(dataSlStr) - 1; i >= 0; i-- {
			res += dataSlStr[i] + "\n"
		}
		ioutil.WriteFile("D:/Golang/L2/Task3/res.txt", []byte(res), 0777)
	case "-u":
		for i, val := range dataSlStr {
			if strings.Count(string(data), val) >= 2 {
				copy(dataSlStr[i:], dataSlStr[i+1:])
			}
		}
		for _, val := range dataSlStr {
			res += val
		}
		ioutil.WriteFile("D:/Golang/L2/Task3/res.txt", []byte(res), 0777)
	default:
		log.Fatal("Unknown command!")
	}

}
