package resources

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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
