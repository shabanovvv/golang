package logger

import (
	"fmt"
	"log"
	"time"
)

const (
	debug = "DEBUG"
	info  = "INFO"
	warn  = "WARN"
	err   = "ERROR"
)

var logLevelPriority = map[string]int{
	debug: 0,
	info:  1,
	warn:  2,
	err:   3,
}

type Logger struct {
	level string
}

func New(level string) *Logger {
	return &Logger{
		level: level,
	}
}

func (l Logger) Debug(msg string) {
	l.Log(msg, debug)
}

func (l Logger) Info(msg string) {
	l.Log(msg, info)
}

func (l Logger) Warn(msg string) {
	l.Log(msg, warn)
}

func (l Logger) Error(msg string) {
	l.Log(msg, err)
}

func (l Logger) Log(msg string, level string) {
	if logLevelPriority[l.level] <= logLevelPriority[level] {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logMessage := fmt.Sprintf("[%s] [%s]: %s", timestamp, level, msg)
		log.Println(logMessage)
	}
}
