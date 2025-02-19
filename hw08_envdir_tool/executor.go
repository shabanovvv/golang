package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, v := range env {
		if v.NeedRemove {
			_ = os.Unsetenv(k)
		} else {
			_ = os.Setenv(k, v.Value)
		}
	}

	allowedCommands := map[string]bool{
		"ls":        true,
		"echo":      true,
		"/bin/bash": true,
		"bash":      true,
		"env":       true,
	}

	// Проверка на наличие команды и её разрешение
	if len(cmd) == 0 || !allowedCommands[cmd[0]] {
		return 1
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Env = os.Environ()
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin

	err := command.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return 1
	}

	return command.ProcessState.ExitCode()
}
