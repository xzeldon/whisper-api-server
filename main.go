package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/xzeldon/whisper-api-server/internal/api"
	"github.com/xzeldon/whisper-api-server/internal/resources"
)

func main() {

	e := echo.New()
	e.HideBanner = true

	args, errParsing := resources.ParseFlags()
	if errParsing != nil {
		e.Logger.Error("Error parsing flags: ", errParsing)
		return
	}

	e.Use(middleware.CORS())

	exePath, errs := os.Executable()
	if errs != nil {
		e.Logger.Error(errs)
		return
	}

	exeDir := filepath.Dir(exePath)

	// Change the working directory to the executable directory
	errs = os.Chdir(exeDir)
	if errs != nil {
		e.Logger.Error(errs)
		return
	}

	cwd, _ := os.Getwd()
	fmt.Println("Current working directory:", cwd)

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	whisperState, err := api.InitializeWhisperState(args.ModelPath, args.Language)

	if err != nil {
		e.Logger.Error(err)
	}

	e.POST("/v1/audio/transcriptions", func(c echo.Context) error {

		return api.Transcribe(c, whisperState)
	})

	e.Logger.Fatal(e.Start("127.0.0.1:3000"))
}
