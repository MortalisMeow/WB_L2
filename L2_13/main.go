package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Flags struct {
	Fields    string // номера полей
	Delimiter string // разделитель полей
	Separated bool   // выводить только строки с разделителем
}

type FieldRange struct {
	Start int
	End   int
}

// Разбирает строку с номерами полей в слайс
func parseFields(fieldsStr string) ([]FieldRange, error) {
	if fieldsStr == "" {
		return nil, fmt.Errorf("не указаны поля для вывода")
	}

	var result []FieldRange
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("некорректный диапазон: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("некорректное число в диапазоне: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("некорректное число в диапазоне: %s", rangeParts[1])
			}

			if start < 1 || end < 1 {
				return nil, fmt.Errorf("номера полей должны быть положительными: %s", part)
			}

			if start > end {
				start, end = end, start
			}

			result = append(result, FieldRange{Start: start, End: end})
		} else {
			field, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("некорректный номер поля: %s", part)
			}

			if field < 1 {
				return nil, fmt.Errorf("номер поля должен быть положительным: %d", field)
			}

			result = append(result, FieldRange{Start: field, End: -1})
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("не указано ни одного поля")
	}

	return result, nil
}

// Проверяет, нужно ли выводить поле с номером
func fieldSelected(fieldNum int, ranges []FieldRange) bool {
	for _, r := range ranges {
		if r.End == -1 {
			// Одиночное поле
			if fieldNum == r.Start {
				return true
			}
		} else {
			// Диапазон
			if fieldNum >= r.Start && fieldNum <= r.End {
				return true
			}
		}
	}
	return false
}

// Обрабатывает аргументы командной строки
func parseFlags() *Flags {
	flg := &Flags{}

	flag.StringVar(&flg.Fields, "f", "", "выбрать поля (колонки)")
	flag.StringVar(&flg.Delimiter, "d", "\t", "использовать другой разделитель")
	flag.BoolVar(&flg.Separated, "s", false, "только строки с разделителем")

	flag.Parse()

	// Проверяем, что указаны поля
	if flg.Fields == "" {
		fmt.Fprintln(os.Stderr, "Ошибка: необходимо указать флаг -f с номерами полей")
		flag.Usage()
		os.Exit(1)
	}

	return flg
}

// Основная логка
func Cut(reader io.Reader, writer io.Writer, flg *Flags) error {
	// Парсим номера полей
	fieldRanges, err := parseFields(flg.Fields)
	if err != nil {
		return fmt.Errorf("ошибка разбора полей: %w", err)
	}

	scanner := bufio.NewScanner(reader)
	delimiter := flg.Delimiter

	for scanner.Scan() {
		line := scanner.Text()
		processLine(line, writer, delimiter, fieldRanges, flg.Separated)
	}

	return scanner.Err()
}

// Обрабатывает одну строку
func processLine(line string, writer io.Writer, delimiter string,
	fieldRanges []FieldRange, separatedOnly bool) {

	// Проверяем, содержит ли строка разделитель
	hasDelimiter := strings.Contains(line, delimiter)

	// Если флаг s установлен и разделителя нет, то пропускаем строку
	if separatedOnly && !hasDelimiter {
		return
	}

	// Если разделителя нет и не установлен s, выводим всю строку
	if !hasDelimiter {
		fmt.Fprintln(writer, line)
		return
	}

	// Разбиваем строку на поля
	fields := strings.Split(line, delimiter)

	// Собираем выбранные поля
	var selectedFields []string
	for i, field := range fields {
		fieldNum := i + 1
		if fieldSelected(fieldNum, fieldRanges) {
			selectedFields = append(selectedFields, field)
		}
	}

	// Выводим только если есть выбранные поля
	if len(selectedFields) > 0 {
		fmt.Fprintln(writer, strings.Join(selectedFields, delimiter))
	}
}

func getInput() io.Reader {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return os.Stdin
	}

	args := flag.Args()
	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка открытия файла '%s': %v\n", args[0], err)
			os.Exit(1)
		}
		return file
	}

	fmt.Fprintln(os.Stderr, "Ошибка: нет входных данных")
	os.Exit(1)
	return nil
}

func main() {
	flg := parseFlags()
	reader := getInput()

	if err := Cut(reader, os.Stdout, flg); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}
