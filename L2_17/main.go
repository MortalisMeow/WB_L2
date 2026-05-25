package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	host := flag.String("h", "", "localhost")
	port := flag.String("p", "23", "port")

	timeout := flag.Duration(
		"timeout",
		10*time.Second,
		"connection timeout",
	)

	// Парсим
	flag.Parse()

	// Аргументы после флагов
	args := flag.Args()

	if len(args) < 2 {
		fmt.Println("usage: go run main.go <host> <port> [--timeout=5s]")
		os.Exit(1)
	}

	address := net.JoinHostPort(*host, *port)

	conn, err := net.DialTimeout(
		"tcp",
		address,
		*timeout,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", address)

	// Канал для сигнализации завершения

	done := make(chan struct{})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		conn.Close()
		close(done)
	}()

	// Горутина для чтения из сокета и вывода в STDOUT
	go func() {
		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(os.Stdout)

		for {
			// Читаем строку из сокета
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					// Сервер закрыл соединение
					fmt.Fprintf(os.Stderr, "\nConnection closed by server\n")
				} else {
					fmt.Fprintf(os.Stderr, "\nRead error: %v\n", err)
				}
				close(done)
				return
			}

			// Выводим полученные данные
			_, err = writer.Write(line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Write to stdout error: %v\n", err)
				close(done)
				return
			}
			writer.Flush()
		}
	}()

	// Горутина для чтения из STDIN и отправки в сокет
	go func() {
		reader := bufio.NewReader(os.Stdin)
		writer := bufio.NewWriter(conn)

		for {
			// Читаем строку из STDIN
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					// Ctrl+D - пользователь хочет завершить соединение
					fmt.Fprintf(os.Stderr, "\nClosing connection...\n")
				} else {
					fmt.Fprintf(os.Stderr, "\nRead from stdin error: %v\n", err)
				}
				close(done)
				return
			}

			// Отправляем данные в сокет
			_, err = writer.Write(line)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Write to socket error: %v\n", err)
				close(done)
				return
			}
			writer.Flush()
		}
	}()

	// Ожидание завершения
	<-done

	// Даем время на корректное закрытие
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Disconnected")

}
