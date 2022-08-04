package main

import (
	"fmt"
	"strconv"
	"strings"
)

/*
=== Задача на распаковку ===
Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)
В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.
Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func convert(s string) string {
	// con := strings.ToLower(str)
	var (
		res string
		err error
	)
	n := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		n[i], err = strconv.Atoi(string(s[i]))
		if err != nil {
			res = res + string(s[i])
		} else {
			res = res + strings.Repeat(string(s[i-1]), n[i]-1)
		}
	}
	// for i, j := 0, 1; i < len(s)-1; i++ {
	// 	n[j], _ = strconv.Atoi(string(s[j]))
	// 	if n[j] != 0 {
	// 		res = res + strings.Repeat(string(s[i]), n[j])
	// 	}
	// 	j++
	// }
	return res

}

func checkout(str string) string {
	strB := []byte(str)
	num := []int{}
	for i := len(strB); i >= 0; i-- {
		if strB[i] == 92 {
			strB = strB[:len(strB)-1]
			continue
		}
		break
	}
	if len(strB) > 0 {

		for i := 0; i < len(strB)-1; i++ {
			if strB[i] == 92 && strB[i+1] == 92 {
				num = append(num, i)
			}
		}
		fmt.Println(num)
		for val := range num {
			copy(strB[:val], strB[val+1:])
			strB = strB[:len(strB)-1]
			// fmt.Println(string(strB))
		}

	}
	return string(strB)
}

func main() {
	var str string
	resStr := ""
	fmt.Scanln(&str)
	// fmt.Println(checkout(str))
	switch {
	case str == "":
		break
	default:
		resStr = convert(str)
	}

	fmt.Println(resStr)
}
