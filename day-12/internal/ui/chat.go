package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/deblinux123/go-ai-rag/day-12/internal/ai"
)

type ChatView struct {
	app *App

	messagesContainer *fyne.Container
	input             *widget.Entry
	sendButton        *widget.Button
}

func NewChatView(app *App) fyne.CanvasObject {
	view := &ChatView{
		app: app,
	}

	title := widget.NewLabelWithStyle(
		"AI Chat",
		fyne.TextAlignLeading,
		fyne.TextStyle{
			Bold: true,
		},
	)

	view.messagesContainer = container.NewVBox()

	scroll := container.NewVScroll(
		view.messagesContainer,
	)

	view.input = widget.NewMultiLineEntry()
	view.input.SetPlaceHolder("Type your message...")

	view.sendButton = widget.NewButton(
		"Send",
		func() {
			view.sendMessage()
		},
	)

	view.input.OnSubmitted = func(_ string) {
		view.sendMessage()
	}

	inputArea := container.NewBorder(
		nil,
		nil,
		nil,
		view.sendButton,
		view.input,
	)

	return container.NewBorder(
		title,
		inputArea,
		nil,
		nil,
		scroll,
	)
}

func (v *ChatView) sendMessage() {
	text := strings.TrimSpace(
		v.input.Text,
	)

	if text == "" {
		return
	}

	if v.app.selectedChatID == 0 {
		v.app.ShowError(
			fmt.Errorf("Please create or select a chat first."),
		)
		return
	}

	_, err := v.app.chatService.CreateMessage(
		v.app.selectedChatID,
		"user",
		text,
	)

	if err != nil {
		v.app.ShowError(err)
		return
	}

	v.addMessage("You", text)

	v.input.SetText("")

	v.sendToOllama(text)
}

func (v *ChatView) sendToOllama(text string) {
	messages := []ai.ChatMessage{
		{
			Role:    "user",
			Content: text,
		},
	}

	tokenChan, errorChan := v.app.aiClient.Chat(
		"qwen2.5:3b",
		messages,
	)

	assistantLabel := widget.NewLabel("")
	assistantLabel.Wrapping = fyne.TextWrapWord

	v.messagesContainer.Add(
		container.NewPadded(
			container.NewVBox(
				widget.NewLabelWithStyle(
					"AI",
					fyne.TextAlignLeading,
					fyne.TextStyle{
						Bold: true,
					},
				),
				assistantLabel,
			),
		),
	)

	go func() {
		var response strings.Builder

		for token := range tokenChan {
			response.WriteString(token)

			current := response.String()

			fyne.Do(func() {
				assistantLabel.SetText(current)
				v.messagesContainer.Refresh()
			})
		}

		if err := <-errorChan; err != nil {
			fyne.Do(func() {
				assistantLabel.SetText(
					"Error: " + err.Error(),
				)
				v.messagesContainer.Refresh()
			})

			return
		}

		_, err := v.app.chatService.CreateMessage(
			v.app.selectedChatID,
			"assistant",
			response.String(),
		)

		if err != nil {
			fyne.Do(func() {
				v.app.ShowError(err)
			})
		}
	}()
}

func (v *ChatView) addMessage(
	role string,
	content string,
) {
	label := widget.NewLabelWithStyle(
		role,
		fyne.TextAlignLeading,
		fyne.TextStyle{
			Bold: true,
		},
	)

	message := widget.NewLabel(content)
	message.Wrapping = fyne.TextWrapWord

	block := container.NewVBox(
		label,
		message,
		widget.NewSeparator(),
	)

	v.messagesContainer.Add(
		container.NewPadded(block),
	)

	v.messagesContainer.Refresh()
}
