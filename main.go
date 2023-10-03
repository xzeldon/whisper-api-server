package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const MODEL_PATH = "./ggml-medium.bin"

func main() {
	e := echo.New()
	e.HideBanner = true

	if l, ok := e.Logger.(*log.Logger); ok {
		l.SetHeader("${time_rfc3339} ${level}")
	}

	whisperState, err := InitializeWhisperState(MODEL_PATH)
	if err != nil {
		e.Logger.Error(err)
		return
	}

	e.POST("/v1/audio/transcriptions", func(c echo.Context) error {
		return transcribe(c, whisperState)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
