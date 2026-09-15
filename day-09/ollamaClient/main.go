package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/deblinux123/go-ai-rag/day-09/ollamaClient/ollama"
)

type ChatScrean struct {
	messages *fyne.Container
	scroll   *container.Scroll

	input       *widget.Entry
	sendButton  *widget.Button
	modelSelect *widget.Select

	status *widget.Label

	ollama *ollama.OllamaClient
}

func NewChatScrean(ollama *ollama.OllamaClient) *ChatScrean {
	messages := container.NewVBox()

	scroll := container.NewVScroll(messages)

	input := widget.NewMultiLineEntry()
	input.SetPlaceHolder("Ask somthing...")

	modelSelect := widget.NewSelect(
		[]string{},
		nil,
	)

	status := widget.NewLabel("Connecting to Ollama...")

	sendButton := widget.NewButton("Send", nil)

	chat := &ChatScrean{
		messages:    messages,
		scroll:      scroll,
		input:       input,
		sendButton:  sendButton,
		modelSelect: modelSelect,
		status:      status,
		ollama:      ollama,
	}

	sendButton.OnTapped = chat.sendMessage

	return chat
}

func (c *ChatScrean) AddMessage(sender, text string) {
	senderLabel := widget.NewLabel(sender)
	messageLabel := widget.NewLabel(text)

	message := container.NewVBox(
		senderLabel,
		messageLabel,
	)

	c.messages.Add(message)
	c.messages.Refresh()

	c.scroll.ScrollToBottom()
}

func (c *ChatScrean) sendMessage() {
	prompt := c.input.Text
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

	c.AddMessage("You", prompt)
	c.status.SetText("AI is thinking...")

	go func() {
		response, err := c.ollama.Chat(model, prompt)

		if err != nil {
			fyne.Do(func() {
				c.status.SetText("Error:" + err.Error())
				c.sendButton.Enable()
			})

			return
		}

		fyne.Do(func() {
			c.AddMessage("AI", response)
			c.status.SetText("Ready")
			c.sendButton.Enable()
		})
	}()
}

func (c *ChatScrean) loadModels() {
	go func() {
		models, err := c.ollama.GetModels()

		if err != nil {
			fyne.Do(func() {
				c.status.SetText("Ollama error: " + err.Error())
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
					fmt.Sprintf("%d models loaded", len(names)),
				)
			} else {
				c.status.SetText("No Ollama models found")
			}

			c.modelSelect.Refresh()
		})
	}()
}

func (c *ChatScrean) Build() fyne.CanvasObject {
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

	w := a.NewWindow("Go AI Assitant")

	ollama := ollama.NewOllamaClient()

	chat := NewChatScrean(ollama)

	title := widget.NewLabel("Go AI Assitant")

	content := container.NewBorder(
		title,
		nil,
		nil,
		nil,
		chat.Build(),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	chat.loadModels()
	w.SetFixedSize(true)
	w.ShowAndRun()
}
