package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type utilFile struct{}

func Copy(fromPath, toPath string, offset, limit int64) error {
	uf := utilFile{}
	from, to, err := uf.Opens(fromPath, toPath)
	if err != nil {
		return err
	}
	defer func(*os.File, *os.File) {
		_ = from.Close()
		_ = to.Close()
	}(from, to)

	fileInfo, err := from.Stat()
	if err != nil {
		return err
	}

	if err = uf.CheckOffsetAndSeek(from, fileInfo, offset); err != nil {
		return err
	}

	buf := make([]byte, 100)
	bytesCopied := int64(0)
	isBreak := false
	for {
		n, errFrom := from.Read(buf)
		if errFrom == io.EOF {
			isBreak = true
		}
		if errFrom != nil && errFrom != io.EOF {
			return errFrom
		}

		if limit > 0 && limit < bytesCopied+int64(n) {
			n = int(limit - bytesCopied)
		}

		m, errTo := to.Write(buf[:n])
		if errTo != nil {
			return errTo
		}
		bytesCopied += int64(m)

		time.Sleep(10 * time.Millisecond)
		uf.ProgressBar(fileInfo.Size(), bytesCopied)

		if isBreak {
			break
		}
	}

	return nil
}

func (uf utilFile) Opens(fromPath, toPath string) (*os.File, *os.File, error) {
	if !uf.IsSupportedTypeFile(fromPath) {
		return nil, nil, ErrUnsupportedFile
	}

	source, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	destination, err := os.Create(toPath)
	if err != nil {
		_ = source.Close()
		return nil, nil, err
	}
	return source, destination, nil
}

func (uf utilFile) IsSupportedTypeFile(filePath string) bool {
	typeExt := []string{".txt", ".bin"}
	ext := filepath.Ext(filePath)
	for _, t := range typeExt {
		if ext == t {
			return true
		}
	}
	return false
}

func (uf utilFile) CheckOffsetAndSeek(file *os.File, fileInfo os.FileInfo, offset int64) error {
	if fileInfo.Size() < 0 {
		return errors.New("file size is unknown")
	}

	if offset > 0 && fileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	} else if offset > 0 {
		_, err := file.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	return nil
}

// ProgressBar — метод для отображения прогресс-бара.
func (uf utilFile) ProgressBar(total int64, current int64) {
	maxWidth := 30 // максимальная ширина прогресс-бара
	if total <= 0 {
		return // избежать деления на ноль
	}

	// Вычисляем процент выполнения
	percent := float64(current) / float64(total)

	// Вычисляем количество символов для завершенной и оставшейся части
	completedWidth := int(percent * float64(maxWidth))
	if completedWidth > maxWidth {
		completedWidth = maxWidth
	}
	keepWidth := maxWidth - completedWidth

	// Формируем строки завершенной и оставшейся части
	completed, keep := "", ""
	if completedWidth > 0 {
		for i := 0; i < completedWidth-1; i++ {
			completed += "="
		}
	}
	for i := 0; i < keepWidth; i++ {
		keep += "_"
	}

	// Выводим прогресс-бар и процент
	fmt.Printf("\r%d/%d [%s>%s] %.2f%%", current, total, completed, keep, percent*100)
}
