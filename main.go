package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xzeldon/whisper-api-server/internal/api"
	"github.com/xzeldon/whisper-api-server/internal/resources"
)

const (
	defaultModelType      = "ggml-medium.bin"
	defaultWhisperVersion = "1.12.0"
)

func changeWorkingDirectory(e *echo.Echo) {
	exePath, err := os.Executable()
	if err != nil {
		e.Logger.Error("Error getting executable path: ", err)
		return
	}

	exeDir := filepath.Dir(exePath)
	if err := os.Chdir(exeDir); err != nil {
		e.Logger.Error("Error changing working directory: ", err)
		return
	}

	cwd, _ := os.Getwd()
	fmt.Println("Current working directory:", cwd)
}

func main() {
	e := echo.New()
	e.HideBanner = true
	changeWorkingDirectory(e)

	args, err := resources.ParseFlags()
	if err != nil {
		e.Logger.Error("Error parsing flags: ", err)
		return
	}

	if _, err := resources.HandleWhisperDll(defaultWhisperVersion); err != nil {
		e.Logger.Error("Error handling Whisper.dll: ", err)
		return
	}

	if _, err := resources.HandleDefaultModel(defaultModelType); err != nil {
		e.Logger.Error("Error handling model file: ", err)
		return
	}

	e.Use(middleware.CORS())

	whisperState, err := api.InitializeWhisperState(args.ModelPath, args.Language)
	if err != nil {
		e.Logger.Error("Error initializing Whisper state: ", err)
		return
	}

	e.POST("/v1/audio/transcriptions", func(c echo.Context) error {
		return api.Transcribe(c, whisperState)
	})

	address := fmt.Sprintf("127.0.0.1:%d", args.Port)
	e.Logger.Fatal(e.Start(address))
}
