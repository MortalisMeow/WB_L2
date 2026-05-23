package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type Flags struct {
	After      int            // количество строк после совпадения (флаг -A)
	Before     int            // количество строк до совпадения (флаг -B)
	Context    int            // количество строк контекста вокруг совпадения (флаг -C)
	Count      bool           // только подсчет совпадений (флаг -c)
	IgnoreCase bool           // игнорировать регистр (флаг -i)
	Invert     bool           // инвертировать фильтр (флаг -v)
	Fixed      bool           // точное совпадение строки (флаг -F)
	LineNum    bool           // выводить номера строк (флаг -n)
	Pattern    string         // шаблон для поиска
	useRegex   bool           // использовать ли регулярное выражение
	compiledRe *regexp.Regexp // скомпилированное регулярное выражение
}

func parseFlags() *Flags {
	flg := &Flags{}

	flag.IntVar(&flg.After, "A", 0, "печатать N строк после совпадения")
	flag.IntVar(&flg.Before, "B", 0, "печатать N строк до совпадения")
	flag.IntVar(&flg.Context, "C", 0, "печатать N строк контекста (до и после)")
	flag.BoolVar(&flg.Count, "c", false, "вывести только количество совпадений")
	flag.BoolVar(&flg.IgnoreCase, "i", false, "игнорировать регистр")
	flag.BoolVar(&flg.Invert, "v", false, "инвертировать фильтр")
	flag.BoolVar(&flg.Fixed, "F", false, "точное совпадение строки")
	flag.BoolVar(&flg.LineNum, "n", false, "выводить номера строк")

	flag.Parse()

	// Если указан флаг -C, он переопределяет -A и -B
	if flg.Context > 0 {
		flg.After = flg.Context
		flg.Before = flg.Context
	}

	// Получаем паттерн (последний аргумент после флагов)
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: не указан шаблон для поиска")
		os.Exit(1)
	}
	flg.Pattern = args[0]

	// Режим поиска
	flg.useRegex = !flg.Fixed

	// Если используем регулярное выражение
	if flg.useRegex {
		pattern := flg.Pattern
		if flg.IgnoreCase {
			pattern = "(?i)" + pattern
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка компиляции регулярного выражения: %v\n", err)
			os.Exit(1)
		}
		flg.compiledRe = re
	}

	return flg
}

// Проверяет соответствуе условиям поиска
func matchLine(flg *Flags, line string) bool {
	var matched bool

	if flg.useRegex {
		// Поиск по регулярному выражению
		matched = flg.compiledRe.MatchString(line)
	} else {
		// Поиск точной подстроки
		if flg.IgnoreCase {
			matched = strings.Contains(strings.ToLower(line), strings.ToLower(flg.Pattern))
		} else {
			matched = strings.Contains(line, flg.Pattern)
		}
	}

	// Инвертируем результат если нужно
	if flg.Invert {
		return !matched
	}
	return matched
}

// Информация о строке для буфера
type lineInfo struct {
	text string
	num  int
}

func printLine(writer io.Writer, line string, lineNum int, showNum bool) {
	if showNum {
		fmt.Fprintf(writer, "%d:%s\n", lineNum, line)
	} else {
		fmt.Fprintln(writer, line)
	}
}

func getInput() io.Reader {
	args := flag.Args()

	if len(args) > 1 {
		file, err := os.Open(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка открытия файла '%s': %v\n", args[1], err)
			os.Exit(1)
		}
		return file
	}

	// Проверяем, есть ли данные в stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return os.Stdin
	}

	fmt.Fprintln(os.Stderr, "Ошибка: не указан входной файл и stdin пуст")
	fmt.Fprintln(os.Stderr, "Использование: program [ФЛАГИ] ШАБЛОН [ФАЙЛ]")
	os.Exit(1)
	return nil
}

func Grep(reader io.Reader, writer io.Writer, flg *Flags) error {
	scanner := bufio.NewScanner(reader)

	// Буфер для хранения предыдущих строк
	beforeBuffer := make([]lineInfo, 0, flg.Before)

	currentLineNum := 0     // текущий номер строки
	matchCount := 0         // количество совпадений для флага -c
	afterCount := 0         // сколько строк еще нужно вывести после совпадения
	printSeparator := false // нужно ли вывести разделитель между группами

	if flg.Count {
		for scanner.Scan() {
			line := scanner.Text()
			if matchLine(flg, line) {
				matchCount++
			}
		}
		fmt.Fprintln(writer, matchCount)
		return scanner.Err()
	}

	for scanner.Scan() {
		line := scanner.Text()
		currentLineNum++

		matched := matchLine(flg, line)

		if matched {
			if printSeparator && (flg.After > 0 || flg.Before > 0) {
				fmt.Fprintln(writer, "--")
			}

			// Выводим строки из буфера
			for _, info := range beforeBuffer {
				printLine(writer, info.text, info.num, flg.LineNum)
			}

			// Выводим текущую строку
			printLine(writer, line, currentLineNum, flg.LineNum)

			afterCount = flg.After
			printSeparator = true

			// Очищаем буфер
			beforeBuffer = beforeBuffer[:0]

		} else if afterCount > 0 {
			// Выводим строки контекста после совпадения
			printLine(writer, line, currentLineNum, flg.LineNum)
			afterCount--

		} else {
			// Добавляем строку в буфер
			if flg.Before > 0 {
				beforeBuffer = append(beforeBuffer, lineInfo{
					text: line,
					num:  currentLineNum,
				})

				// Удаляем самые старые строки
				if len(beforeBuffer) > flg.Before {
					beforeBuffer = beforeBuffer[1:]
				}
			}
		}
	}

	return scanner.Err()
}

func main() {
	flg := parseFlags()
	reader := getInput()

	if err := Grep(reader, os.Stdout, flg); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при обработке: %v\n", err)
		os.Exit(1)
	}

}
