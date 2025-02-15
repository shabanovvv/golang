package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-envdir /path/to/env/dir command arg1 arg2 ...")
		os.Exit(1)
	}

	envDir := os.Args[1]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Println("Error reading environment directory:", err)
		os.Exit(1)
	}

	// Передаем команду и её аргументы
	command := os.Args[2:]
	exitCode := RunCmd(command, env)

	os.Exit(exitCode)
}
