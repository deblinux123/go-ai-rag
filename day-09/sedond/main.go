package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ChatScrean struct {
	messages   *fyne.Container
	scroll     *container.Scroll
	input      *widget.Entry
	sendButton *widget.Button
}

func NewChatScrean() *ChatScrean {
	messages := container.NewVBox()

	scroll := container.NewVScroll(messages)

	input := widget.NewEntry()
	input.SetPlaceHolder("Ask somthing...")

	sendButton := widget.NewButton("Send", nil)

	chat := &ChatScrean{
		messages:   messages,
		scroll:     scroll,
		input:      input,
		sendButton: sendButton,
	}

	sendButton.OnTapped = chat.sendMessage

	input.OnSubmitted = func(s string) {
		chat.sendMessage()
	}

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
	text := c.input.Text

	if text == "" {
		return
	}

	c.AddMessage("You", text)

	c.input.SetText("")

	c.AddMessage("AI", fmt.Sprintf("This is a face AI response to: %s", text))
}

func (c *ChatScrean) Build() fyne.CanvasObject {
	inputArea := container.NewBorder(
		nil,
		nil,
		nil,
		c.sendButton,
		c.input,
	)

	return container.NewBorder(
		nil,
		inputArea,
		nil,
		nil,
		c.scroll,
	)
}
func main() {
	a := app.New()

	w := a.NewWindow("Go AI Assistant")

	chat := NewChatScrean()

	title := widget.NewLabel("Go AI Assistant")

	content := container.NewBorder(
		title,
		nil,
		nil,
		nil,
		chat.Build(),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
