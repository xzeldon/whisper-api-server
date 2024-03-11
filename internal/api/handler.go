package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type TranscribeResponse struct {
	Text string `json:"text"`
}

func TranscribeFromFile(c echo.Context, whisperState *WhisperState) error {
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

	if err != nil {
		c.Logger().Errorf("Error processing audio: %s", err)
		return err
	}

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

func Transcribe(c echo.Context, whisperState *WhisperState) error {
	// Get the file header
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.Logger().Errorf("Error retrieving the file: %s", err)
		return err
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.Logger().Errorf("Error opening the file: %s", err)
		return err
	}
	defer file.Close()

	// Read the file into a buffer
	buffer, err := io.ReadAll(file)
	if err != nil {
		c.Logger().Errorf("Error reading the file into buffer: %s", err)
		return err
	}

	whisperState.mutex.Lock()
	defer whisperState.mutex.Unlock()

	bufferSpecial, err := whisperState.media.LoadAudioFileData(&buffer, true)

	if err != nil {
		c.Logger().Errorf("Error loading audio file data: %s", err)
		return err
	}

	err = whisperState.context.RunStreamed(whisperState.params, bufferSpecial)
	if err != nil {
		c.Logger().Errorf("Error processing audio: %s", err)
		return err
	}

	result, err := getResult(whisperState.context)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	if len(result) == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	response := TranscribeResponse{
		Text: strings.TrimLeft(result, " "),
	}

	return c.JSON(http.StatusOK, response)
}
