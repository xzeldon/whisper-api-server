# Whisper API Server (Go)

## ⚠️ This project is a work in progress (WIP).

This API server enables audio transcription using the OpenAI Whisper models.

# Setup

- Download the desired model from [huggingface](https://huggingface.co/ggerganov/whisper.cpp/tree/main)
- Update the model path in the `main.go` file
- Download `Whisper.dll` from [github](https://github.com/Const-me/Whisper/releases/tag/1.12.0) (`Library.zip`) and place it in the project's root directory
- Build project: `go build .` (you only need go compiler, without gcc)

# Usage example

Make a request to the server using the following command:

```sh
curl http://localhost:3000/v1/audio/transcriptions \
  -H "Content-Type: multipart/form-data" \
  -F file="@/path/to/file/audio.mp3" \
```

Receive a response in JSON format:

```json
{
  "text": "Imagine the wildest idea that you've ever had, and you're curious about how it might scale to something that's a 100, a 1,000 times bigger. This is a place where you can get to do that."
}
```

# Roadmap

- [x] Implement automatic model downloading from [huggingface](https://huggingface.co/ggerganov/whisper.cpp/tree/main)
- [x] Implement automatic `Whisper.dll` downloading from [Guthub releases](https://github.com/Const-me/Whisper/releases)
- [ ] Provide prebuilt binaries for Windows
- [ ] Include instructions for running on Linux with Wine (likely possible).

# Credits

- [Const-me/Whisper](https://github.com/Const-me/Whisper) project
- [goConstmeWhisper](https://github.com/jaybinks/goConstmeWhisper) for the remarkable Go bindings for [Const-me/Whisper](https://github.com/Const-me/Whisper)
- [Georgi Gerganov](https://github.com/ggerganov) for GGML models
