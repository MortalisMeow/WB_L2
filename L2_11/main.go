package main

import (
	"fmt"
	"sort"
	"strings"
)

// функция для сортировки слова для дальнейшего определения анаграммы
func sortString(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

// убирает дубликаты
func removeDuplicates(strs []string) []string {
	if len(strs) == 0 {
		return strs
	}
	result := []string{strs[0]}
	for i := 1; i < len(strs); i++ {
		if strs[i] != strs[i-1] {
			result = append(result, strs[i])
		}
	}
	return result
}

func findAnagramm(s []string) map[string][]string {
	tempMap := make(map[string][]string, 5)

	//проходим по слайсу и добавляем во временную мапу
	for _, word := range s {
		lower := strings.ToLower(word)

		key := sortString(lower)

		tempMap[key] = append(tempMap[key], lower)
	}

	result := make(map[string][]string)

	for _, v := range tempMap {
		if len(v) < 2 {
			continue
		}
		sort.Strings(v)

		result[v[0]] = removeDuplicates(v)
	}
	return result

}

func main() {
	anagramms := []string{"пятак", "тяпка", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

	mapAnagramms := findAnagramm(anagramms)

	fmt.Println(mapAnagramms)

}
