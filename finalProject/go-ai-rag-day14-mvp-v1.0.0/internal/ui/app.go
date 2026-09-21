package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"

	"github.com/deblinux123/go-ai-rag/day-14/internal/ai"
	"github.com/deblinux123/go-ai-rag/day-14/internal/chat"
	"github.com/deblinux123/go-ai-rag/day-14/internal/settings"
)

type App struct {
	fyneApp fyne.App
	window  fyne.Window

	chatService *chat.Service
	aiClient    *ai.Client
	settings    *settings.Store

	selectedChatID int64

	sidebar      *Sidebar
	chatView     *ChatView
	settingsView *SettingsView
}

func NewApp(
	fyneApp fyne.App,
	chatService *chat.Service,
	aiClient *ai.Client,
	store *settings.Store,
) *App {
	a := &App{
		fyneApp:     fyneApp,
		window:      fyneApp.NewWindow("Go AI RAG"),
		chatService: chatService,
		aiClient:    aiClient,
		settings:    store,
	}

	a.applyTheme()

	return a
}

func (a *App) UpdateAIClientURL(url string) {
	if url == "" {
		return
	}

	a.aiClient.BaseURL = url
}

func (a *App) Run() {
	a.window.Resize(fyne.NewSize(1280, 800))
	a.window.SetMaster()

	a.sidebar = NewSidebar(a)
	a.chatView = NewChatView(a)
	a.settingsView = NewSettingsView(a)

	main := container.NewHSplit(
		a.sidebar.CanvasObject(),
		a.chatView.CanvasObject(),
	)
	main.SetOffset(0.24)

	content := container.NewBorder(
		nil,
		nil,
		nil,
		a.settingsView.CanvasObject(),
		main,
	)

	a.window.SetContent(content)
	a.window.ShowAndRun()
}

func (a *App) applyTheme() {
	cfg := a.settings.Get()

	var variant fyne.ThemeVariant
	if cfg.Theme == "light" {
		variant = theme.VariantLight
	} else {
		variant = theme.VariantDark
	}

	a.fyneApp.Settings().SetTheme(newAppTheme(variant))
}

func (a *App) SetTheme(name string) {
	if name != "light" {
		name = "dark"
	}

	if err := a.settings.Update(func(cfg *settings.Config) {
		cfg.Theme = name
	}); err != nil {
		a.ShowError(err)
		return
	}

	a.applyTheme()
}

func (a *App) Config() settings.Config {
	return a.settings.Get()
}

func (a *App) ShowError(err error) {
	if err != nil {
		dialog.ShowError(err, a.window)
	}
}

func (a *App) ShowInfo(title, message string) {
	dialog.ShowInformation(title, message, a.window)
}

func (a *App) DeleteSelectedChat() {
	if a.chatView != nil && a.chatView.IsBusy() {
		return
	}

	if a.selectedChatID == 0 {
		return
	}

	id := a.selectedChatID

	dialog.ShowConfirm(
		"Delete conversation",
		"Delete this conversation and all of its messages?",
		func(ok bool) {
			if !ok {
				return
			}

			if err := a.chatService.DeleteChat(id); err != nil {
				a.ShowError(err)
				return
			}

			a.selectedChatID = 0
			a.sidebar.Reload()
			a.chatView.Clear()
		},
		a.window,
	)
}

func (a *App) SelectChat(id int64) {
	if a.chatView != nil && a.chatView.IsBusy() {
		return
	}

	a.selectedChatID = id

	if a.chatView != nil {
		a.chatView.LoadChat(id)
	}
}

func (a *App) CreateChat() {
	if a.chatView != nil && a.chatView.IsBusy() {
		return
	}

	id, err := a.chatService.CreateChat("New Chat")
	if err != nil {
		a.ShowError(err)
		return
	}

	a.selectedChatID = id
	a.sidebar.Reload()
	a.sidebar.SelectByID(id)
	a.chatView.LoadChat(id)
}

func (a *App) RefreshModels() {
	if a.chatView != nil {
		a.chatView.SetBusy(true)
	}

	go func() {
		cfg := a.Config()
		a.aiClient.SetBaseURL(cfg.OllamaURL)

		models, err := a.aiClient.GetModels()

		fyne.Do(func() {
			if a.chatView != nil {
				a.chatView.SetBusy(false)
			}

			if err != nil {
				a.ShowError(err)
				return
			}

			a.settingsView.SetModels(models)
		})
	}()
}
