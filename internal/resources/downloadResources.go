package resources

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

func GetModel(modelType string) (string, error) {
	fileURL := fmt.Sprintf("https://huggingface.co/ggerganov/whisper.cpp/resolve/main/%s", modelType)
	filePath := modelType

	isModelFileExists := IsFileExists(filePath)

	if !isModelFileExists {
		fmt.Println("Model not found.")
		err := DownloadFile(fileURL, filePath)
		if err != nil {
			return "", err
		}
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	fmt.Printf("Model found: %s\n", absPath)
	return filePath, nil
}

func DownloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileSize := resp.ContentLength
	bar := progressbar.DefaultBytes(
		fileSize,
		"Downloading",
	)

	writer := io.MultiWriter(out, bar)

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetWhisperDll(version string) (string, error) {
	fileUrl := fmt.Sprintf("https://github.com/Const-me/Whisper/releases/download/%s/Library.zip", version)
	fileToExtract := "Binary/Whisper.dll"

	isWhisperDllExists := IsFileExists("Whisper.dll")

	if !isWhisperDllExists {
		fmt.Println("Whisper DLL not found.")
		archivePath, err := os.CreateTemp("", "WhisperLibrary-*.zip")
		if err != nil {
			return "", err
		}
		defer archivePath.Close()

		err = DownloadFile(fileUrl, archivePath.Name())
		if err != nil {
			return "", err
		}

		err = extractFile(archivePath.Name(), fileToExtract)
		if err != nil {
			return "", err
		}
	}

	absPath, err := filepath.Abs("Whisper.dll")
	if err != nil {
		return "", err
	}

	fmt.Printf("Library found: %s\n", absPath)
	return "Whisper.dll", nil
}

func extractFile(archivePath string, fileToExtract string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		if file.Name == fileToExtract {
			targetPath := filepath.Base(fileToExtract)

			writer, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer writer.Close()

			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			_, err = io.Copy(writer, src)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("File not found in the archive")
}

func IsFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}