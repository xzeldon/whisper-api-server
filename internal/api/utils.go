package api

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func saveFormFile(name string, c echo.Context) (string, error) {
	file, err := c.FormFile(name)
	if err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	tmpDir, err := ensureDir("tmp")
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	filename := time.Now().Format(time.RFC3339)
	filename = tmpDir + "/" + sanitizeFilename(filename) + ext

	dst, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filename, nil
}

func sanitizeFilename(filename string) string {
	invalidChars := []string{`\`, `/`, `:`, `*`, `?`, `"`, `<`, `>`, `|`}
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, char, "-")
	}
	return filename
}

func ensureDir(dirPath string) (string, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0700)
		if err != nil {
			return "", err
		}
	}

	return dirPath, nil
}
