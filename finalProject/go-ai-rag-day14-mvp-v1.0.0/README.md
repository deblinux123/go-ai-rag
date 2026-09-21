# Go AI RAG

A modern local-first desktop AI chat client built with Go, Fyne, SQLite and Ollama.

## Features

### AI
- Ollama integration
- Dynamic model selection
- Streaming responses
- System prompt
- Temperature control

### Conversations
- New conversations
- Persistent chat history
- Multiple conversations
- Delete conversations
- Scrollable message history
- Automatic conversation titles

### Settings
- Ollama URL
- Model
- Temperature
- Dark / light theme
- System prompt
- Test Ollama connection

### Storage
- SQLite
- Persistent conversations
- Per-user application data directory

## Architecture

```text
Fyne UI
   |
   v
Chat Service
   |
   +------------------+
   |                  |
   v                  v
Ollama Client      SQLite
   |
   v
Ollama
   |
   v
Local LLM
```

## Project structure

```text
day-14/
├── cmd/
│   └── app/
│       ├── main.go
│       ├── FyneApp.toml
│       └── Icon.png
├── internal/
│   ├── ai/
│   │   └── ollama.go
│   ├── chat/
│   │   ├── models.go
│   │   └── service.go
│   ├── database/
│   │   └── sqlite.go
│   ├── settings/
│   │   └── settings.go
│   └── ui/
│       ├── app.go
│       ├── chat.go
│       ├── settings.go
│       ├── sidebar.go
│       └── theme.go
├── .github/
│   └── workflows/
│       └── release.yml
├── build-linux.sh
├── go.mod
├── go.sum
└── README.md
```

## Requirements

- Go 1.27.1+
- Ollama
- A local Ollama model
- GCC and Linux graphics development packages for building Fyne

On Ubuntu/Debian:

```bash
sudo apt update
sudo apt install -y gcc libgl1-mesa-dev xorg-dev libwayland-dev libxkbcommon-dev
```

Fyne's Linux build prerequisites include Go, GCC and the OpenGL/X11/Wayland development headers.

## Run

From this directory:

```bash
go mod tidy
go run ./cmd/app
```

Make sure Ollama is running:

```bash
ollama serve
```

And pull a model:

```bash
ollama pull qwen2.5:3b
```

The default Ollama URL is:

```text
http://localhost:11434
```

## Build Linux binary

```bash
chmod +x build-linux.sh
./build-linux.sh
```

Or directly:

```bash
go build   -trimpath   -ldflags="-s -w"   -o dist/go-ai-rag   ./cmd/app
```

Run it:

```bash
./dist/go-ai-rag
```

## Desktop package

Install the Fyne CLI:

```bash
go install fyne.io/tools/cmd/fyne@latest
```

Then:

```bash
cd cmd/app
fyne package --os linux --release
```

Fyne creates a Linux archive containing the desktop application structure.

## Data location

The application stores its data in the operating system's user configuration directory.

It creates:

```text
GoAIRAG/
├── settings.json
└── conversations.db
```

On Linux this is normally under:

```text
~/.config/GoAIRAG/
```

## Versioning

This project uses semantic version tags.

Create the first release:

```bash
git add .
git commit -m "feat: release v1.0.0"

git tag -a v1.0.0 -m "Go AI RAG v1.0.0"

git push origin main
git push origin v1.0.0
```

Pushing a `v1.0.0` tag triggers GitHub Actions and publishes the Linux binary and desktop package to the GitHub Release.

For the next release:

```bash
git add .
git commit -m "feat: improve chat experience"

git tag -a v1.1.0 -m "Go AI RAG v1.1.0"

git push origin main
git push origin v1.1.0
```

## Release types

- `v1.0.0` — first stable MVP
- `v1.0.1` — bug fix
- `v1.1.0` — backward-compatible feature
- `v2.0.0` — breaking architectural or product change

## Security / privacy

This application is local-first. Chat messages are stored locally in SQLite and requests are sent to the configured Ollama server.

No cloud AI API is required by the application itself.
