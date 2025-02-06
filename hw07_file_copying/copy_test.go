package main

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	fileSource, err := os.Open("testdata/input.txt")
	require.NoError(t, err)
	tempFile, err := os.CreateTemp("testdata", "test_file_*.txt")
	require.NoError(t, err)
	/*defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("error removing temp file %s\n", name)
		}
	}(tempFile.Name())*/

	t.Run("offset: 0, limit: 10", func(t *testing.T) {
		var offset int64
		var limit int64 = 10

		err = Copy("testdata/input.txt", tempFile.Name(), offset, limit)
		require.NoError(t, err)

		_, err := fileSource.Seek(offset, io.SeekStart)
		require.NoError(t, err)
		buf := make([]byte, limit)
		bytesRead, err := fileSource.Read(buf)
		require.NoError(t, err)
		expectedContent := buf[:bytesRead]

		copiedContent, err := os.ReadFile(tempFile.Name())
		require.NoError(t, err)

		require.Equal(t, expectedContent, copiedContent)
	})

	t.Run("offset: 100, limit: 1000", func(t *testing.T) {
		var offset int64 = 100
		var limit int64 = 1000

		err = Copy("testdata/input.txt", tempFile.Name(), offset, limit)
		require.NoError(t, err)

		_, err := fileSource.Seek(offset, io.SeekStart)
		require.NoError(t, err)
		buf := make([]byte, limit)
		bytesRead, err := fileSource.Read(buf)
		require.NoError(t, err)
		expectedContent := buf[:bytesRead]

		copiedContent, err := os.ReadFile(tempFile.Name())
		require.NoError(t, err)

		require.Equal(t, expectedContent, copiedContent)
	})

	t.Run("ErrUnsupportedFile", func(t *testing.T) {
		err = Copy("testdata/input.mp3", tempFile.Name(), 0, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})
}
