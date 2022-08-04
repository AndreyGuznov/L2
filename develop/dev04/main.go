package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===
Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.
Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	hash      = make(map[string]struct{})
	resultMap = make(map[string][]string)
)

func getLower(sl []string) []string {
	str := ""
	for i := 0; i < len(sl); i++ {
		sl[i] = strings.ToLower(sl[i])
		str += sl[i]
		if strings.Count(str, sl[i]) > 1 {
			sl[i] = ""
		}
	}
	return sl
}

func sortString(word string) string {
	w := strings.Split(word, "")
	sort.Strings(w)
	return strings.Join(w, "")
}

func Anagrams(word string, words []string) []string {
	word = sortString(word)
	var results []string

	for _, w := range words {
		if word == sortString(w) {

			results = append(results, w)

		}
	}

	return results
}

func inputToHash(sl []string) {
	for _, val := range sl {
		hash[val] = struct{}{}
	}
	for key := range hash {
		anagram := Anagrams(key, sl)
		if len(anagram) > 1 {
			resultMap[key] = append(resultMap[key], anagram...)
		}
		for i := 0; i < len(anagram); i++ {
			delete(hash, anagram[i])
		}
	}
}

func main() {
	mass := [...]string{"пяТак", "слиток", "пятка", "стол", "слиток", "стул", "Тяпка", "листок"}
	sl := getLower(mass[:])
	// inputToHash(sl)
	// fmt.Println(&resultMap)
	// mass[0] = "стилок"
	inputToHash(sl)
	fmt.Println(resultMap)
}
