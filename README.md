# Whisper API Server (Go)

## ⚠️ This project is a work in progress (WIP).

This API server enables audio transcription using the OpenAI Whisper models.

# Setup

- Download `.exe` from [Releases](https://github.com/xzeldon/whisper-api-server/releases/latest)
- Just run it!

# Build from source (Windows)

## Prerequisites

- GCC Compiler Installed in your PATH (You can get it from [here](https://github.com/niXman/mingw-builds-binaries))
- Install Go (https://go.dev/doc/install)

Before build make sure that **CGO_ENABLED** env is set to **1**

```
$env:CGO_ENABLED = "1"
```

you can check this with this command

```
go env
```

Also you have to have installed gcc x64 i.e. by MYSYS

Download the sources and use `go build`.
For example, you can build using the following command:

```bash
go build -ldflags "-s -w" -o server.exe main.go
```

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

# Usage with [Obsidian](https://obsidian.md/)

1. Install [Obsidian voice recognotion plugin](https://github.com/nikdanilov/whisper-obsidian-plugin)
2. Open the plugin's settings.
3. Set the following values:
   - API KEY: `sk-1`
   - API URL: `http://localhost:3000/v1/audio/transcriptions`
   - Model: `whisper-1`

# Roadmap

- [x] Implement automatic model downloading from [huggingface](https://huggingface.co/ggerganov/whisper.cpp/tree/main)
- [x] Implement automatic `Whisper.dll` downloading from [Guthub releases](https://github.com/Const-me/Whisper/releases)
- [x] Provide prebuilt binaries for Windows
- [ ] Include instructions for running on Linux with Wine (likely possible).
- [x] Use flags to override the model path
- [x] Use flags to override the port

# Credits

- [Const-me/Whisper](https://github.com/Const-me/Whisper) project
- [goConstmeWhisper](https://github.com/jaybinks/goConstmeWhisper) for the remarkable Go bindings for [Const-me/Whisper](https://github.com/Const-me/Whisper)
- [Georgi Gerganov](https://github.com/ggerganov) for GGML models
