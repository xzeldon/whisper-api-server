package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/xzeldon/whisper-api-server/internal/api"
	"github.com/xzeldon/whisper-api-server/internal/resources"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	_, err := resources.GetWhisperDll("1.12.0")
	if err != nil {
		e.Logger.Error(err)
	}

	model, err := resources.GetModel("ggml-medium.bin")
	if err != nil {
		e.Logger.Error(err)
	}

	whisperState, err := api.InitializeWhisperState(model)
	if err != nil {
		e.Logger.Error(err)
	}

	e.POST("/v1/audio/transcriptions", func(c echo.Context) error {
		return api.Transcribe(c, whisperState)
	})

	e.Logger.Fatal(e.Start("127.0.0.1:3000"))
}
