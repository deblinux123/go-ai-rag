package ui

import (
	"fmt"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/deblinux123/go-ai-rag/day-14/internal/ai"
	"github.com/deblinux123/go-ai-rag/day-14/internal/settings"
)

type ChatView struct {
	app *App

	canvas fyne.CanvasObject

	titleLabel  *widget.Label
	statusLabel *widget.Label

	messages *fyne.Container
	scroll   *container.Scroll

	input *widget.Entry
	send  *widget.Button

	busyMu sync.Mutex
	busy   bool
}

func NewChatView(app *App) *ChatView {
	v := &ChatView{app: app}

	v.titleLabel = widget.NewLabelWithStyle(
		"New conversation",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	v.statusLabel = widget.NewLabel("Ready")

	v.messages = container.NewVBox()
	v.scroll = container.NewVScroll(v.messages)

	v.input = widget.NewMultiLineEntry()
	v.input.SetMinRowsVisible(3)
	v.input.SetPlaceHolder("Message Ollama…")
	v.input.Wrapping = fyne.TextWrapWord

	v.send = widget.NewButtonWithIcon(
		"Send",
		theme.ConfirmIcon(),
		func() {
			v.Send()
		},
	)
	v.send.Importance = widget.HighImportance

	v.input.OnSubmitted = func(string) {
		v.Send()
	}

	clearButton := widget.NewButtonWithIcon(
		"",
		theme.ContentClearIcon(),
		func() {
			v.input.SetText("")
		},
	)
	clearButton.Importance = widget.LowImportance

	header := container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewVBox(
			v.titleLabel,
			v.statusLabel,
		),
	)

	composer := container.NewBorder(
		nil,
		nil,
		nil,
		v.send,
		container.NewBorder(
			nil,
			nil,
			nil,
			clearButton,
			v.input,
		),
	)

	v.canvas = container.NewBorder(
		header,
		container.NewPadded(composer),
		nil,
		nil,
		v.scroll,
	)

	v.Clear()

	return v
}

func (v *ChatView) CanvasObject() fyne.CanvasObject {
	return v.canvas
}

func (v *ChatView) Clear() {
	v.messages.Objects = nil
	v.messages.Refresh()

	v.titleLabel.SetText("New conversation")
	v.statusLabel.SetText("Ready")
	v.input.SetText("")
}

func (v *ChatView) SetBusy(value bool) {
	v.busyMu.Lock()
	v.busy = value
	v.busyMu.Unlock()

	if value {
		v.send.Disable()
		v.statusLabel.SetText("Generating…")
	} else {
		v.send.Enable()
		v.statusLabel.SetText("Ready")
	}
}

func (v *ChatView) IsBusy() bool {
	v.busyMu.Lock()
	defer v.busyMu.Unlock()
	return v.busy
}

func (v *ChatView) LoadChat(id int64) {
	messages, err := v.app.chatService.GetMessages(id)
	if err != nil {
		v.app.ShowError(err)
		return
	}

	chatInfo, err := v.app.chatService.GetChat(id)
	if err != nil {
		v.app.ShowError(err)
		return
	}

	v.messages.Objects = nil

	for _, message := range messages {
		v.addMessage(message.Role, message.Content)
	}

	v.titleLabel.SetText(chatInfo.Title)
	v.messages.Refresh()
	v.scroll.ScrollToBottom()
	v.input.SetText("")
}

func (v *ChatView) addMessage(role, content string) *widget.Label {
	title := "AI"
	if role == "user" {
		title = "You"
	}

	body := widget.NewLabel(content)
	body.Wrapping = fyne.TextWrapWord
	body.Selectable = true

	card := widget.NewCard(
		title,
		"",
		container.NewPadded(body),
	)

	v.messages.Add(card)
	v.messages.Refresh()
	v.scroll.ScrollToBottom()

	return body
}

func (v *ChatView) Send() {
	if v.IsBusy() {
		return
	}

	if v.app.selectedChatID == 0 {
		v.app.ShowError(fmt.Errorf("please create or select a chat first"))
		return
	}

	text := strings.TrimSpace(v.input.Text)
	if text == "" {
		return
	}

	v.input.SetText("")
	v.addMessage("user", text)

	cfg := v.app.Config()

	if _, err := v.app.chatService.CreateMessage(
		v.app.selectedChatID,
		"user",
		text,
	); err != nil {
		v.app.ShowError(err)
		return
	}

	if chatInfo, err := v.app.chatService.GetChat(v.app.selectedChatID); err == nil && chatInfo.Title == "New Chat" {
		title := makeTitle(text)
		_ = v.app.chatService.UpdateChatTitle(v.app.selectedChatID, title)
		v.titleLabel.SetText(title)
		v.app.sidebar.Reload()
	}

	v.SetBusy(true)
	go v.streamResponse(cfg)
}

func (v *ChatView) streamResponse(cfg settings.Config) {
	v.app.aiClient.SetBaseURL(cfg.OllamaURL)

	chatID := v.app.selectedChatID

	history, err := v.app.chatService.GetMessages(chatID)
	if err != nil {
		fyne.Do(func() {
			v.SetBusy(false)
			v.app.ShowError(err)
		})
		return
	}

	messages := make([]ai.ChatMessage, 0, len(history)+1)

	if strings.TrimSpace(cfg.SystemPrompt) != "" {
		messages = append(messages, ai.ChatMessage{
			Role:    "system",
			Content: cfg.SystemPrompt,
		})
	}

	for _, message := range history {
		messages = append(messages, ai.ChatMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	tokens, errors := v.app.aiClient.Chat(
		cfg.Model,
		messages,
		cfg.Temperature,
	)

	var response strings.Builder
	var assistantLabel *widget.Label
	tokenOpen := true
	errorOpen := true

	for tokenOpen || errorOpen {
		select {
		case token, ok := <-tokens:
			if !ok {
				tokenOpen = false
				continue
			}

			response.WriteString(token)
			current := response.String()

			fyne.Do(func() {
				if assistantLabel == nil {
					assistantLabel = v.addMessage("assistant", "")
				}

				assistantLabel.SetText(current)
				v.messages.Refresh()
				v.scroll.ScrollToBottom()
			})

		case err, ok := <-errors:
			if !ok {
				errorOpen = false
				continue
			}

			if err != nil {
				fyne.Do(func() {
					v.SetBusy(false)
					v.app.ShowError(err)
				})
				return
			}

			errorOpen = false
		}
	}

	if response.Len() == 0 {
		fyne.Do(func() {
			v.SetBusy(false)
			v.app.ShowError(fmt.Errorf("ollama returned an empty response"))
		})
		return
	}

	if _, err := v.app.chatService.CreateMessage(
		chatID,
		"assistant",
		response.String(),
	); err != nil {
		fyne.Do(func() {
			v.SetBusy(false)
			v.app.ShowError(err)
		})
		return
	}

	fyne.Do(func() {
		v.SetBusy(false)
		v.app.sidebar.Reload()
	})
}

func (v *ChatView) SetBusyForModelRefresh(value bool) {
	v.SetBusy(value)
}

func makeTitle(text string) string {
	text = strings.Join(strings.Fields(text), " ")

	if len([]rune(text)) <= 36 {
		return text
	}

	runes := []rune(text)
	return string(runes[:36]) + "…"
}
