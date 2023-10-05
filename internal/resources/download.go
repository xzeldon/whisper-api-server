package resources

import (
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

func DownloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileSize := resp.ContentLength
	bar := progressbar.DefaultBytes(
		fileSize,
		"Downloading",
	)

	writer := io.MultiWriter(out, bar)

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
