package resources

import (
	"fmt"
	"path/filepath"
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
