package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xzeldon/whisper-api-server/internal/api"
	"github.com/xzeldon/whisper-api-server/internal/resources"
)

// begin delimiter const
const beginDelimiter = "[begin]"
const endDelimiter = "[end]"

func change_working_directory() {
	exePath, errs := os.Executable()
	if errs != nil {
		println("Error getting executable path")
		return
	}

	exeDir := filepath.Dir(exePath)

	// Change the working directory to the executable directory
	errs = os.Chdir(exeDir)
	if errs != nil {
		println("Error changing working directory")
		return
	}

	cwd, _ := os.Getwd()
	fmt.Println("Current working directory:", cwd)
}

func main() {

	change_working_directory()

	args, errParsing := resources.ParseFlags()
	if errParsing != nil {
		println("Error parsing flags: ", errParsing)
		return
	}

	whisperState, err := api.InitializeWhisperState(args.ModelPath, args.Language)

	if err != nil {
		println("Error initializing whisper state: ", err)
	}
	const maxCapacity = 2048 * 10240

	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	println("waiting_for_input")
	if scanner.Scan() {
		base64Data := scanner.Text()
		decodedBuffer, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			fmt.Println("Error decoding buffer:", err)
			return
		}

		result := api.TranscribeBytes(decodedBuffer, whisperState)
		println(beginDelimiter + result + endDelimiter)
		println("finished")

		// Process the decodedBuffer (e.g., print its length)
		fmt.Println("Received buffer size:", len(decodedBuffer))

		// Send a response back to Node.js (optional)
		fmt.Fprintln(os.Stdout, "Buffer received successfully!")
	} else if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
	}

	// e.Logger.Fatal(e.Start(fmt.Sprintf("127.0.0.1:%d", args.Port)))
}
