package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Обработка аргументов командной строки
	timeout := flag.Duration("timeout", 10*time.Second, "timeout for connecting to the server")
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Println("Usage: go-telnet [--timeout=duration] host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := fmt.Sprintf("%s:%s", host, port)

	// Создание клиента Telnet
	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	// Подключение к серверу
	if err := client.Connect(); err != nil {
		fmt.Printf("Error connecting to %s: %v\n", address, err)
		os.Exit(1)
	}
	defer client.Close()

	// Обработка сигналов для завершения программы
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	// Запуск горутины для приема данных от сервера
	go func() {
		if err := client.Receive(); err != nil {
			fmt.Printf("Error reading from server: %v\n", err)
		}
		fmt.Println("...Connection was closed by peer")
		os.Exit(0)
	}()

	// Основной цикл для отправки сообщений от STDIN
	for {
		select {
		case <-sigs:
			fmt.Println("\n...Exiting")
			return
		default:
			if err := client.Send(); err != nil {
				fmt.Printf("Error sending message: %v\n", err)
				return
			}
		}
	}
}
