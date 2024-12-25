package resources

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PromptUser prompts the user with a question and returns true if they agree
func PromptUser(question string) bool {
	fmt.Printf("%s (y/n): ", question)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// HandleWhisperDll checks if Whisper.dll exists or prompts the user to download it
func HandleWhisperDll(version string) (string, error) {
	if IsFileExists("Whisper.dll") {
		absPath, err := filepath.Abs("Whisper.dll")
		if err != nil {
			return "", err
		}
		fmt.Printf("Library found: %s\n", absPath)
		return "Whisper.dll", nil
	}

	fmt.Println("Whisper DLL not found.")
	if PromptUser("Do you want to download Whisper.dll automatically?") {
		path, err := GetWhisperDll(version)
		if err != nil {
			return "", fmt.Errorf("failed to download Whisper.dll: %w", err)
		}
		return path, nil
	}

	fmt.Println("To use Whisper, download the DLL manually:")
	fmt.Printf("URL: https://github.com/Const-me/Whisper/releases/download/%s/Library.zip\n", version)
	fmt.Println("Extract 'Binary/Whisper.dll' from the archive and place it in the executable's directory.")
	fmt.Println("You can manually specify path to .dll file using cli arguments, use --help to print available cli flags")
	return "", fmt.Errorf("whisper.dll not found and user chose not to download")
}

// HandleDefaultModel checks if the default model exists or prompts the user to download it
func HandleDefaultModel(modelType string) (string, error) {
	if IsFileExists(modelType) {
		absPath, err := filepath.Abs(modelType)
		if err != nil {
			return "", err
		}
		fmt.Printf("Model found: %s\n", absPath)
		return modelType, nil
	}

	fmt.Println("Default model not found.")
	if PromptUser("Do you want to download the default model (ggml-medium.bin) automatically?") {
		path, err := GetModel(modelType)
		if err != nil {
			return "", fmt.Errorf("failed to download the default model: %w", err)
		}
		return path, nil
	}

	fmt.Println("To use Whisper, download the model manually:")
	fmt.Println("URL: https://huggingface.co/ggerganov/whisper.cpp/tree/main")
	fmt.Println("Place the model file in the executable's directory or specify its path using cli arguments.")
	fmt.Println("You can manually specify path to model file using cli arguments, use --help to print available cli flags")
	return "", fmt.Errorf("default model not found and user chose not to download")
}
