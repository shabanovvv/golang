package logger

import (
	"bytes"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf) // Перенаправляем вывод в буфер

	tests := []struct {
		level     string
		logLevel  string
		message   string
		shouldLog bool
	}{
		{debug, debug, "Debug message", true},
		{debug, info, "Info message", true},
		{debug, warn, "Warning message", true},
		{debug, err, "Error message", true},
		{info, debug, "Debug message", false},
		{info, info, "Info message", true},
		{info, warn, "Warning message", true},
		{info, err, "Error message", true},
		{warn, debug, "Debug message", false},
		{warn, info, "Info message", false},
		{warn, warn, "Warning message", true},
		{warn, err, "Error message", true},
		{err, debug, "Debug message", false},
		{err, info, "Info message", false},
		{err, warn, "Warning message", false},
		{err, err, "Error message", true},
	}

	for _, tt := range tests {
		buf.Reset() // Очищаем буфер перед каждым тестом
		logger := New(tt.level)
		logger.Log(tt.message, tt.logLevel)

		if tt.shouldLog && buf.Len() == 0 {
			t.Errorf("expected log message for level %s, got none", tt.level)
		} else if !tt.shouldLog && buf.Len() > 0 {
			t.Errorf("did not expect log message for level %s, but got: %s", tt.level, buf.String())
		}
	}
}

func TestLoggerInitialization(t *testing.T) {
	logger := New(warn)
	if logger.level != warn {
		t.Errorf("expected log level to be %s, got %s", warn, logger.level)
	}
}
