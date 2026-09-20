package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/deblinux123/go-ai-rag/day-12/internal/ai"
	"github.com/deblinux123/go-ai-rag/day-12/internal/chat"
)

type App struct {
	fyneApp fyne.App
	window  fyne.Window

	chatService *chat.Service
	aiClient    *ai.Client

	selectedChatID int64
}

func NewAPP(fyneApp fyne.App, chatService *chat.Service, aiClient *ai.Client) *App {
	return &App{
		fyneApp:     fyneApp,
		window:      fyneApp.NewWindow("Go AI Desktop"),
		chatService: chatService,
		aiClient:    aiClient,
	}
}

func (a *App) Run() {
	a.window.Resize(fyne.NewSize(1100, 700))

	sidebar := NewSidebar(a)
	chatView := NewChatView(a)
	settings := NewSettingsView(a)

	content := container.NewHSplit(
		sidebar,
		container.NewBorder(
			nil,
			nil,
			nil,
			settings,
			chatView,
		),
	)

	content.SetOffset(0.25)

	a.window.SetContent(content)

	a.window.ShowAndRun()
}
