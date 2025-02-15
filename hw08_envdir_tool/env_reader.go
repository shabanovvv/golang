package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var ErrNotExistEnvDir = errors.New("environment directory does not exist")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, ErrNotExistEnvDir
	}

	for _, file := range files {
		if !file.IsDir() {
			firstLine, err := GetFirstLine(dir + string(os.PathSeparator) + file.Name())
			if err != nil {
				return nil, err
			}

			env[file.Name()] = EnvValue{
				strings.TrimRight(firstLine, " \t"),
				len(firstLine) == 0,
			}
		}
	}
	return env, nil
}

func GetFirstLine(pathFile string) (string, error) {
	file, err := os.Open(pathFile)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "\x00", "\n")
		return line, nil
	} else if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}
