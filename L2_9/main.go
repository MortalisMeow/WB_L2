package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var (
	ErrInvalidString = errors.New("некорректная строка")
)

func StringUnpack(input string) (string, error) {
	//обработка пустой строки
	if len(input) == 0 {
		return "", nil
	}

	isTrue := false
	var runes = []rune(input)

	//начинается с цифры
	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	//буфер для конкатенации
	var sb strings.Builder

	for i := 0; i < len(runes); i++ {
		v := runes[i]
		if unicode.IsLetter(v) {
			sb.WriteRune(v)
			//есть хотя бы одна цифра
			isTrue = true
			if i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
				//повторяем букву на -1 кол-во, т.к уже добавили ранее
				sb.WriteString(strings.Repeat(string(v), int((runes[i+1]-'0')-1)))
				i++
			}
		}
	}
	if isTrue == false {
		return "", ErrInvalidString
	}
	return sb.String(), nil

}

func main() {
	var input string
	var output string

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите строку: ")

	if scanner.Scan() {
		input = scanner.Text()
	}

	output, err := StringUnpack(input)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(output)
	}

}
