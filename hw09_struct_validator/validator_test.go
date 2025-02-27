package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "Valid User",
			in: User{
				ID:     strings.Repeat("a", 36),
				Name:   "Test User",
				Age:    30,
				Email:  "test@example.com",
				Role:   UserRole("admin, stuff"),
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			name: "Invalid User ID Length",
			in: User{
				ID:     "too-short",
				Name:   "Test User",
				Age:    30,
				Email:  "test@example.com",
				Role:   UserRole("admin, stuff"),
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationError{Field: "ID", Err: errors.New("field ID is invalid")},
		},
		{
			name: "Invalid User Age",
			in: User{
				ID:     strings.Repeat("a", 36),
				Name:   "Test User",
				Age:    10,
				Email:  "test@example.com",
				Role:   UserRole("admin, stuff"),
				Phones: []string{"12345678901"},
			},
			expectedErr: errors.New("field 'Age': field Age is invalid"),
		},
		{
			name: "Invalid User Email",
			in: User{
				ID:     strings.Repeat("a", 36),
				Name:   "Test User",
				Age:    30,
				Email:  "invalid-email",
				Role:   UserRole("admin, stuff"),
				Phones: []string{"12345678901"},
			},
			expectedErr: errors.New("field 'Email': field Email is invalid"),
		},
		{
			name: "Invalid User Role",
			in: User{
				ID:     strings.Repeat("a", 36),
				Name:   "Test User",
				Age:    30,
				Email:  "test@example.com",
				Role:   "invalid-role",
				Phones: []string{"12345678901"},
			},
			expectedErr: errors.New("field 'Role': field Role is invalid"),
		},
		{
			name: "Invalid User Phones Length",
			in: User{
				ID:     strings.Repeat("a", 36),
				Name:   "Test User",
				Age:    30,
				Email:  "test@example.com",
				Role:   UserRole("admin, stuff"),
				Phones: []string{"123"},
			},
			expectedErr: errors.New("field 'Phones': field Phones is invalid"),
		},
		{
			name: "Valid App",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid App Version Length",
			in: App{
				Version: "1.0",
			},
			expectedErr: errors.New("field 'Version': field Version is invalid"),
		},
		{
			name: "Valid Response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: errors.New("invalid tag format; expected 'in:<value1,value2,...>', got 'in:200,404,500'"),
		},
		{
			name: "Invalid Response Code",
			in: Response{
				Code: 300,
				Body: "OK",
			},
			expectedErr: errors.New("invalid tag format; expected 'in:<value1,value2,...>', got 'in:200,404,500'"),
		},
		{
			name: "Invalid User Details",
			in: User{
				ID:     "123456",               // некорректное значение
				Name:   "",                     // пустое имя
				Age:    15,                     // возраст менее 18
				Email:  "invalid-email-format", // некорректный email
				Role:   UserRole("guest"),      // некорректная роль
				Phones: []string{"short"},      // некорректный номер телефона
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: errors.New("field ID is invalid")},
				{Field: "Age", Err: errors.New("field Age is invalid")},
				{Field: "Email", Err: errors.New("field Email is invalid")},
				{Field: "Role", Err: errors.New("field Role is invalid")},
				{Field: "Phones", Err: errors.New("field Phones is invalid")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)

				var validationErrors ValidationErrors
				if !errors.As(tt.expectedErr, &validationErrors) {
					assert.EqualError(t, err, tt.expectedErr.Error(), "Сообщение об ошибке не соответствует ожидаемому")
				}
			}
		})
	}
}
