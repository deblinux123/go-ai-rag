# 🚀 Learn Fyne with Go — 14 Day Roadmap

A practical 14-day roadmap for learning **Fyne GUI framework with Golang** by building a real AI desktop application.

The final goal is to build a desktop AI application using:

* 🐹 Go
* 🎨 Fyne
* 🤖 Ollama
* 🌐 REST API
* 🗄️ SQLite
* 🧠 AI / LLM
* ⚡ Streaming responses

---

## 🎯 Final Project

During these 14 days, we will build an AI desktop application:

```text
┌──────────────────────────────────────────────┐
│              🤖 Go AI Assistant              │
├───────────────┬──────────────────────────────┤
│               │                              │
│  Chats        │       AI Conversation        │
│               │                              │
│  + New Chat   │  You: Explain RAG            │
│               │                              │
│  Chat 1       │  AI: RAG stands for...       │
│  Chat 2       │                              │
│  Chat 3       │                              │
│               │                              │
├───────────────┴──────────────────────────────┤
│  Ask something...                  [Send]    │
└──────────────────────────────────────────────┘
```

The application will communicate with Ollama:

```text
Fyne UI
   │
   ▼
Go Application
   │
   ▼
Ollama API
   │
   ▼
LLM
   │
   ▼
Streaming Response
   │
   ▼
Fyne UI
```

---

# 📚 Prerequisites

Before starting, you should know the basics of:

* Go syntax
* Variables
* Functions
* Structs
* Interfaces
* Goroutines
* Channels
* HTTP requests
* JSON
* Basic REST APIs
* Git/GitHub

You don't need to be an expert.

---

# 🗓️ 14-Day Roadmap

## Day 01 — Fyne Fundamentals

### 🎯 Goal

Understand what Fyne is and create your first desktop application.

### Learn

* What is Fyne?
* Fyne architecture
* Application
* Window
* Canvas
* Widgets
* Containers
* Layouts
* `app.New()`
* `app.NewWithID()`
* `window.ShowAndRun()`

### Practice

Create:

```text
Hello Fyne
```

Then create:

```text
Hello AI
```

with a button.

### Mini Project

Build:

```text
Hello AI Desktop
```

Features:

* Window
* Label
* Button
* Button click event

---

# Day 02 — Widgets

### 🎯 Goal

Learn the most important Fyne widgets.

### Learn

* Label
* Button
* Entry
* MultiLineEntry
* Select
* Check
* Radio
* SelectEntry
* ProgressBar
* Slider

### Practice

Create a form:

```text
Name:
[____________]

Model:
[ qwen2.5:3b ▼ ]

Temperature:
[------●-----]

[ Start AI ]
```

### Mini Project

Create:

```text
AI Settings Form
```

---

# Day 03 — Containers & Layouts

### 🎯 Goal

Learn how to create real application layouts.

### Learn

* `container.NewVBox`
* `container.NewHBox`
* `container.NewGridWithColumns`
* `container.NewBorder`
* `container.NewCenter`
* `container.NewStack`

### Practice

Build:

```text
┌──────────────────────────────┐
│          Header              │
├───────────┬──────────────────┤
│           │                  │
│ Sidebar   │ Main Content     │
│           │                  │
│           │                  │
├───────────┴──────────────────┤
│          Footer              │
└──────────────────────────────┘
```

### Mini Project

Create the skeleton of the final AI application.

---

# Day 04 — Events & User Interaction

### 🎯 Goal

Understand event-driven programming in Fyne.

### Learn

* Button callbacks
* Entry events
* Select events
* Keyboard interaction
* Updating widgets
* Form validation

### Practice

Create a login-like form without authentication:

```text
Username
Password

[Submit]
```

When the user clicks submit:

```text
Welcome Farid!
```

### Mini Project

Create:

```text
AI Configuration Screen
```

with:

* Model
* Temperature
* System prompt
* Save button

---

# Day 05 — Fyne Data Binding

### 🎯 Goal

Learn how to manage application state.

### Learn

* Data binding
* Binding strings
* Binding booleans
* Binding values
* Observing changes
* Updating UI from state

### Practice

Create:

```text
AI State

Model: qwen2.5:3b
Temperature: 0.7
Streaming: true
```

Change the UI and synchronize the internal state.

### Mini Project

Create an application state:

```go
type AppState struct {
    Model       string
    Temperature float64
    Streaming   bool
}
```

---

# Day 06 — Goroutines & Fyne UI

### 🎯 Goal

Learn how to run background tasks without freezing the UI.

This is extremely important for the AI application.

### Learn

* Goroutines
* Channels
* Background workers
* UI updates
* Thread-safe UI patterns
* Loading indicators

### Practice

Create:

```text
[ Start ]

Processing...
██████████░░░░░░░
```

Run the work in a goroutine.

### Mini Project

Create a fake AI response generator:

```text
Thinking...

Hello! I am your AI assistant...
```

---

# Day 07 — HTTP Client with Go

### 🎯 Goal

Connect the Fyne application to an external API.

### Learn

* `net/http`
* HTTP GET
* HTTP POST
* Request headers
* JSON encoding
* JSON decoding
* HTTP errors
* Timeouts

### Practice

Create:

```go
type Client struct {
    BaseURL string
}
```

and implement:

```go
func (c *Client) Get(...)
func (c *Client) Post(...)
```

### Mini Project

Create an API client package:

```text
internal/
└── api/
    └── client.go
```

---

# Day 08 — Ollama Integration

### 🎯 Goal

Connect your Fyne application to Ollama.

### Architecture

```text
Fyne
  │
  ▼
Go Ollama Client
  │
  ▼
localhost:11434
  │
  ▼
Ollama
  │
  ▼
LLM
```

### Learn

* Ollama API
* Chat endpoint
* Model endpoint
* JSON request
* JSON response

### Create

```go
type OllamaClient struct {
    BaseURL string
}
```

Implement:

```go
Chat()
ListModels()
ShowModel()
```

### Mini Project

Send:

```text
Hello
```

to Ollama and display:

```text
Hello! How can I help you?
```

inside Fyne.

---

# Day 09 — Build Chat UI

### 🎯 Goal

Build the actual AI chat interface.

### UI

```text
┌─────────────────────────────────────┐
│          Go AI Assistant             │
├─────────────────────────────────────┤
│                                     │
│ You                                 │
│ Explain RAG                         │
│                                     │
│ AI                                  │
│ RAG is a technique...               │
│                                     │
├─────────────────────────────────────┤
│ Ask something...                    │
│                                     │
│                         [ Send ]     │
└─────────────────────────────────────┘
```

### Learn

* Scroll containers
* Chat messages
* Dynamic widgets
* Adding widgets at runtime
* Refreshing containers

### Mini Project

Build:

```text
ChatScreen
```

---

# Day 10 — Streaming AI Responses

### 🎯 Goal

Make the AI response appear token-by-token.

Instead of:

```text
[wait 10 seconds]

Hello, this is the complete response.
```

we want:

```text
Hello
Hello, this
Hello, this is
Hello, this is the AI
Hello, this is the AI response
```

### Learn

* HTTP streaming
* JSON streams
* `bufio.Scanner`
* Channels
* Goroutines
* UI updates

### Architecture

```text
Ollama
  │
  │ token
  ▼
Go Channel
  │
  ▼
Fyne UI
  │
  ▼
Chat Message
```

### Mini Project

Implement:

```go
func StreamChat(...)
```

---

# Day 11 — Chat History

### 🎯 Goal

Store conversations.

### Learn

* Chat models
* Conversation state
* SQLite
* CRUD
* Persistence

### Data model

```go
type Chat struct {
    ID        int
    Title     string
    CreatedAt time.Time
}
```

```go
type Message struct {
    ID        int
    ChatID    int
    Role      string
    Content   string
    CreatedAt time.Time
}
```

### Database

```text
SQLite
   │
   ├── chats
   │
   └── messages
```

### Mini Project

Implement:

```text
New Chat
Delete Chat
Load Chat
Save Message
Load Messages
```

---

# Day 12 — Application Architecture

### 🎯 Goal

Refactor the project into a professional structure.

Recommended structure:

```text
go-ai-desktop/
│
├── cmd/
│   └── app/
│       └── main.go
│
├── internal/
│   │
│   ├── ai/
│   │   └── ollama.go
│   │
│   ├── chat/
│   │   ├── service.go
│   │   └── models.go
│   │
│   ├── database/
│   │   └── sqlite.go
│   │
│   └── ui/
│       ├── app.go
│       ├── chat.go
│       ├── sidebar.go
│       └── settings.go
│
├── assets/
│
├── go.mod
├── go.sum
└── README.md
```

### Learn

* Separation of concerns
* Services
* Interfaces
* Dependency injection
* UI layer
* AI layer
* Database layer

### Goal

Avoid:

```text
main.go = 2000 lines
```

Instead:

```text
UI
 ↓
Service
 ↓
AI Client
 ↓
Ollama
```

---

# Day 13 — Settings & UX

### 🎯 Goal

Make the application feel like a real product.

### Add

* Model selection
* Temperature
* System prompt
* Dark/light theme
* Ollama URL
* Clear conversation
* New conversation
* Delete conversation
* Loading indicator
* Error messages

### Settings

```text
┌─────────────────────────────┐
│ Settings                    │
├─────────────────────────────┤
│ Ollama URL                  │
│ [ http://localhost:11434 ]  │
│                             │
│ Model                       │
│ [ qwen2.5:3b ▼ ]            │
│                             │
│ Temperature                 │
│ [-------●------]             │
│                             │
│ System Prompt               │
│ [ You are a helpful AI... ] │
│                             │
│             [ Save ]        │
└─────────────────────────────┘
```

---

# Day 14 — Final Product

### 🎯 Goal

Finish the MVP.

Your application should support:

### 🤖 AI

* Ollama
* Model selection
* Chat
* Streaming
* System prompt

### 💬 Chat

* New chat
* Chat history
* Delete chat
* Multiple conversations
* Scrollable messages

### ⚙️ Settings

* Ollama URL
* Model
* Temperature
* Theme
* System prompt

### 💾 Storage

* SQLite
* Persistent conversations

### 🖥️ Desktop

* Linux
* Windows
* macOS

---

# 🏗️ Final Architecture

```text
                 ┌───────────────┐
                 │   Fyne UI     │
                 └───────┬───────┘
                         │
                         ▼
                 ┌───────────────┐
                 │ Chat Service  │
                 └───────┬───────┘
                         │
              ┌──────────┴──────────┐
              │                     │
              ▼                     ▼
      ┌───────────────┐     ┌───────────────┐
      │ Ollama Client │     │ SQLite DB     │
      └───────┬───────┘     └───────────────┘
              │
              ▼
      ┌───────────────┐
      │    Ollama     │
      └───────┬───────┘
              │
              ▼
          ┌───────┐
          │  LLM  │
          └───────┘
```

---

# 📅 Daily Workflow

Every day follow this process:

```text
30 min → Learn
60 min → Code
30 min → Build a mini project
20 min → Debug
10 min → Write README / notes
```

Total:

```text
≈ 2.5 hours/day
```

If you have more time:

```text
3–4 hours/day
```

---

# 📊 Progress Tracker

## Week 1

* [ ] Day 01 — Fyne Fundamentals
* [ ] Day 02 — Widgets
* [ ] Day 03 — Containers & Layouts
* [ ] Day 04 — Events
* [ ] Day 05 — Data Binding
* [ ] Day 06 — Goroutines & UI
* [ ] Day 07 — HTTP Client

## Week 2

* [ ] Day 08 — Ollama Integration
* [ ] Day 09 — Chat UI
* [ ] Day 10 — Streaming
* [ ] Day 11 — SQLite & History
* [ ] Day 12 — Architecture
* [ ] Day 13 — Settings & UX
* [ ] Day 14 — Final MVP

---

# 🧪 Final Feature Checklist

```text
[ ] Fyne application
[ ] Multiple screens
[ ] Sidebar
[ ] Chat UI
[ ] Ollama connection
[ ] Model selector
[ ] Chat completion
[ ] Streaming response
[ ] SQLite
[ ] Chat history
[ ] New chat
[ ] Delete chat
[ ] Settings
[ ] Dark mode
[ ] Error handling
[ ] Loading state
[ ] Clean architecture
```

---

# 🚀 Future Roadmap

After completing these 14 days, continue with:

## Phase 2 — AI Features

```text
RAG
Embeddings
Qdrant
Document upload
PDF processing
Semantic search
Context retrieval
```

Architecture:

```text
             ┌─────────────┐
             │   Fyne UI   │
             └──────┬──────┘
                    │
                    ▼
              Go Backend
                    │
          ┌─────────┼─────────┐
          ▼         ▼         ▼
       Ollama     Qdrant    SQLite
          │         │
          ▼         ▼
        LLM      Embeddings
```

---

# Phase 3 — Product Features

Add:

* 📄 PDF upload
* 📝 TXT/Markdown support
* 🔎 RAG search
* 🧠 Memory
* 📚 Knowledge bases
* 💬 Multiple AI models
* 🎤 Speech-to-text
* 🔊 Text-to-speech
* 🖼️ Vision models
* 🔌 API integrations

---

# Phase 4 — Production

Learn:

* Application packaging
* Cross-platform builds
* Auto updates
* Logging
* Configuration management
* Crash reporting
* Performance optimization
* Security
* Licensing

---

# 🎯 Final Goal

The goal isn't simply:

> "I learned Fyne."

The goal is:

> **I can build and ship a real AI desktop application using Go.**

Technology stack:

```text
             AI Desktop Product
                    │
        ┌───────────┴───────────┐
        │                       │
      Fyne                     Go
        │                       │
        └───────────┬───────────┘
                    │
                 Ollama
                    │
                   LLM
                    │
              ┌─────┴─────┐
              │           │
            Qdrant      SQLite
```

---

# ⭐ Learning Philosophy

Don't spend 14 days only watching tutorials.

For every concept:

```text
Learn
  ↓
Code
  ↓
Break it
  ↓
Debug it
  ↓
Build something
  ↓
Commit to Git
```

Recommended Git workflow:

```bash
git add .
git commit -m "day-01: learn fyne fundamentals"
git push
```

At the end:

```bash
git log --oneline
```

You should have approximately:

```text
14 days
14+ commits
14+ mini projects
1 complete AI desktop MVP
```

---

# 🏁 Definition of Done

You have successfully completed this roadmap when you can run:

```bash
go run ./cmd/app
```

and get a desktop application where you can:

```text
1. Open the application
2. Select an Ollama model
3. Create a new conversation
4. Send a message
5. Receive a streaming response
6. See the conversation in the UI
7. Close the application
8. Reopen it
9. Continue the previous conversation
```



