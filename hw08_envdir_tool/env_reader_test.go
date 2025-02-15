package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("directory does not exist", func(t *testing.T) {
		_, err := ReadDir("./321123312/")
		require.Truef(t, errors.Is(err, ErrNotExistEnvDir), "%v", err)
	})

	t.Run("directory exists and check env", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)
		require.Len(t, env, 5)

		expectedEnv := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}
		require.Equal(t, expectedEnv, env)
	})
}
