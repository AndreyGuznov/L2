package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
=== Взаимодействие с ОС ===
Необходимо реализовать собственный шелл
встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах
Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		executeCommand(strings.Fields(scanner.Text()))
	}
}
func executeCommand(commands []string) {
	args := commands[1:]
	command := commands[0]

	switch command {
	case "quit":
		fmt.Println("Bye, bye")
		os.Exit(1)
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(dir)
	case "ps":
		c := exec.Command("TASKLIST")
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	case "cd":
		dir := strings.Join(args, "")
		os.Chdir(dir)
	case "echo":
		fmt.Println(strings.Join(args, " "))
	case "kill":
		pid, err := strconv.Atoi(strings.Join(args, ""))
		if err != nil {
			fmt.Println(1)
			fmt.Println(err)
			break
		}
		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Println(err)
			break
		}
		err = proc.Kill()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Process killed")
	default:
		fmt.Println("Unknown command")
	}
}
