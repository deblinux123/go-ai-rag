package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/deblinux123/go-ai-rag/day-10/Client/ollama"
)

type ChatScreen struct {
	messages *fyne.Container
	scroll   *container.Scroll

	input       *widget.Entry
	sendButton  *widget.Button
	modelSelect *widget.Select
	status      *widget.Label

	ollama *ollama.OllamaClient
}

func NewChatScreen(ollama *ollama.OllamaClient) *ChatScreen {
	messages := container.NewVBox()

	scroll := container.NewVScroll(messages)

	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Ask something...")

	modelSelect := widget.NewSelect(
		[]string{},
		nil,
	)

	status := widget.NewLabel("Connecting to Ollama...")

	sendButton := widget.NewButton("Send", nil)

	chat := &ChatScreen{
		messages:    messages,
		scroll:      scroll,
		input:       input,
		sendButton:  sendButton,
		modelSelect: modelSelect,
		status:      status,
		ollama:      ollama,
	}

	sendButton.OnTapped = chat.sendMessage

	input.OnSubmitted = func(_ string) {
		chat.sendMessage()
	}

	return chat
}

func (c *ChatScreen) AddMessage(sender string, text string) *widget.Label {
	senderLabel := widget.NewLabel(sender)

	messageLabel := widget.NewLabel(text)
	messageLabel.Wrapping = fyne.TextWrapWord

	message := container.NewVBox(
		senderLabel,
		messageLabel,
	)

	c.messages.Add(message)
	c.messages.Refresh()

	c.scroll.ScrollToBottom()

	return messageLabel
}

func (c *ChatScreen) sendMessage() {
	prompt := strings.TrimSpace(c.input.Text)
	model := c.modelSelect.Selected

	if prompt == "" {
		return
	}

	if model == "" {
		c.status.SetText("Please select a model")
		return
	}

	c.input.SetText("")
	c.sendButton.Disable()
	c.status.SetText("AI is thinking...")

	// Add user's message.
	c.AddMessage("You", prompt)

	// Create an empty AI message.
	aiMessage := c.AddMessage("AI", "")

	tokens, errors := c.ollama.StreamChat(
		model,
		prompt,
	)

	go func() {
		var response strings.Builder

		for {
			select {
			case token, ok := <-tokens:
				if !ok {
					fyne.Do(func() {
						c.status.SetText("Ready")
						c.sendButton.Enable()
					})
					return
				}

				response.WriteString(token)

				currentText := response.String()

				fyne.Do(func() {
					aiMessage.SetText(currentText)
					c.messages.Refresh()
					c.scroll.ScrollToBottom()
				})

			case err, ok := <-errors:
				if !ok {
					continue
				}

				fyne.Do(func() {
					aiMessage.SetText(
						"Error: " + err.Error(),
					)

					c.messages.Refresh()
					c.status.SetText("Error")
					c.sendButton.Enable()
				})

				return
			}
		}
	}()
}

func (c *ChatScreen) loadModels() {
	go func() {
		models, err := c.ollama.GetModels()

		if err != nil {
			fyne.Do(func() {
				c.status.SetText(
					"Ollama error: " + err.Error(),
				)
			})
			return
		}

		names := make([]string, 0, len(models))

		for _, model := range models {
			names = append(names, model.Name)
		}

		fyne.Do(func() {
			c.modelSelect.Options = names

			if len(names) > 0 {
				c.modelSelect.SetSelected(names[0])

				c.status.SetText(
					fmt.Sprintf(
						"%d models loaded",
						len(names),
					),
				)
			} else {
				c.status.SetText(
					"No Ollama models found",
				)
			}

			c.modelSelect.Refresh()
		})
	}()
}

func (c *ChatScreen) Build() fyne.CanvasObject {
	inputArea := container.NewBorder(
		nil,
		nil,
		nil,
		c.sendButton,
		c.input,
	)

	top := container.NewVBox(
		c.modelSelect,
		c.status,
	)

	return container.NewBorder(
		top,
		inputArea,
		nil,
		nil,
		c.scroll,
	)
}

func main() {
	a := app.New()

	w := a.NewWindow("Go AI Assistant")

	ollama := ollama.NewOllamaClient()

	chat := NewChatScreen(ollama)

	title := widget.NewLabel("Go AI Assistant")

	content := container.NewBorder(
		title,
		nil,
		nil,
		nil,
		chat.Build(),
	)

	w.SetContent(content)

	w.Resize(fyne.NewSize(700, 700))

	chat.loadModels()

	w.ShowAndRun()
}
