package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// Флаги
	column := flag.Int("k", 0, "номер колонки для сортировки (начиная с 1)")
	numeric := flag.Bool("n", false, "сортировать по числовому значению")
	reverse := flag.Bool("r", false, "сортировать в обратном порядке")
	unique := flag.Bool("u", false, "не выводить повторяющиеся строки")

	flag.Parse()

	// Читаем строки
	lines := readLines()

	// Сортируем
	sortLines(lines, *column, *numeric, *reverse)

	// Убираем дубликаты если нужно
	if *unique {
		lines = removeDuplicates(lines)
	}

	// Выводим
	for _, line := range lines {
		fmt.Println(line)
	}
}

func readLines() []string {
	var scanner *bufio.Scanner

	// Если передан файл
	if flag.NArg() > 0 {
		filename := flag.Arg(0)
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка открытия файла: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		// Иначе читаем из STDIN
		scanner = bufio.NewScanner(os.Stdin)
	}

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
		os.Exit(1)
	}

	return lines
}

// sortLines сортирует строки с учётом флагов
func sortLines(lines []string, column int, numeric bool, reverse bool) {
	sort.Slice(lines, func(i, j int) bool {
		// Получаем значения для сравнения
		a := getValue(lines[i], column)
		b := getValue(lines[j], column)

		var less bool

		if numeric {
			// Числовая сортировка
			less = compareNumeric(a, b)
		} else {
			// Строковая сортировка
			less = a < b
		}

		// Инвертируем
		if reverse {
			return !less
		}
		return less
	})
}

// getValue возвращает значение из нужной колонки или всю строку
func getValue(line string, column int) string {
	if column <= 0 {
		return line
	}

	// Делим по табуляции
	parts := strings.Split(line, "\t")

	index := column - 1

	if index < len(parts) {
		return parts[index]
	}

	return ""
}

// compareNumeric сравнивает строки как числа
func compareNumeric(a, b string) bool {
	// Пробуем преобразовать в числа
	numA, errA := strconv.ParseFloat(strings.TrimSpace(a), 64)
	numB, errB := strconv.ParseFloat(strings.TrimSpace(b), 64)

	if errA == nil && errB == nil {
		return numA < numB
	}

	return a < b
}

// removeDuplicates убирает повторяющиеся строки
func removeDuplicates(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	result := []string{lines[0]}

	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			result = append(result, lines[i])
		}
	}

	return result
}
