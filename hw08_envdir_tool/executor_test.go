package main

import (
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name         string
		cmd          []string
		env          Environment
		shouldFail   bool
		expectedCode int
	}{
		{
			name:         "Set variable and run command",
			cmd:          []string{"env"},
			env:          Environment{"TEST_VAR": {Value: "TEST_VALUE", NeedRemove: false}},
			shouldFail:   false,
			expectedCode: 0,
		},
		{
			name:         "Unset variable and run command",
			cmd:          []string{"env"},
			env:          Environment{"TEST_VAR": {Value: "", NeedRemove: true}},
			shouldFail:   false,
			expectedCode: 0,
		},
		{
			name:         "Run invalid command",
			cmd:          []string{"invalid_command"},
			env:          Environment{},
			shouldFail:   true,
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Установка переменных окружения перед выполнением команды
			returnCode := RunCmd(tt.cmd, tt.env)

			if tt.shouldFail && returnCode == 0 {
				t.Errorf("expected failure, got success")
			}

			if !tt.shouldFail && returnCode != 0 {
				t.Errorf("expected success, got failure with code: %d", returnCode)
			}

			// Проверка значений переменных окружения
			for k, v := range tt.env {
				if v.NeedRemove {
					if _, exists := os.LookupEnv(k); exists {
						t.Errorf("expected environment variable %s to be unset", k)
					}
				} else {
					val, exists := os.LookupEnv(k)
					if !exists || val != v.Value {
						t.Errorf("expected environment variable %s to be %s, got %s", k, v.Value, val)
					}
				}
			}
		})
	}
}
