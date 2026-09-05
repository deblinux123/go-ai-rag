package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("Hello AI")

	messageLabel := widget.NewLabel("Message:")

	message := widget.NewEntry()
	message.SetPlaceHolder("Ask something...")

	responseLabel := widget.NewLabel("Response:")

	response := widget.NewLabel("")

	button := widget.NewButton("Ask", func() {
		response.SetText("You asked: " + message.Text)
	})

	responseBox := container.NewVBox(
		responseLabel,
		response,
	)

	content := container.NewVBox(
		messageLabel,
		message,
		button,
		responseBox,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
