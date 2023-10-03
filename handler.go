package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type TranscribeResponse struct {
	Text string `json:"text"`
}

func transcribe(c echo.Context, whisperState *WhisperState) error {
	audioPath, err := saveFormFile("file", c)
	if err != nil {
		c.Logger().Errorf("Error reading file: %s", err)
		return err
	}

	whisperState.mutex.Lock()
	buffer, err := whisperState.media.LoadAudioFile(audioPath, true)
	if err != nil {
		c.Logger().Errorf("Error loading audio file data: %s", err)
	}

	err = whisperState.context.RunFull(whisperState.params, buffer)

	result, err := getResult(whisperState.context)
	if err != nil {
		c.Logger().Error(err)
	}

	defer whisperState.mutex.Unlock()

	if len(result) == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := TranscribeResponse{
		Text: strings.TrimLeft(result, " "),
	}

	return c.JSON(http.StatusOK, response)
}
