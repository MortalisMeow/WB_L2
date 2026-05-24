package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	//Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			fmt.Println()
		}
	}()

	for {
		fmt.Print("myshell> ")

		if !scanner.Scan() {
			fmt.Println("\nexit")
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" || line == "exit" {
			if line == "exit" {
				return
			}
			continue
		}

		//Обработка команд
		args := strings.Fields(line)
		switch args[0] {
		case "cd":
			if len(args) > 1 {
				os.Chdir(args[1])
			}
		case "pwd":
			if dir, err := os.Getwd(); err == nil {
				fmt.Println(dir)
			}
		default:
			executeLine(line)
		}
	}
}

func executeLine(line string) {
	var cmd *exec.Cmd

	if strings.Contains(line, "|") {
		cmd = exec.Command("sh", "-c", line)
	} else {
		args := strings.Fields(line)
		cmd = exec.Command(args[0], args[1:]...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
