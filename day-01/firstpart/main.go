package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("com.farid.helloai")

	label := widget.NewLabel("Hello AI")

	button := widget.NewButton("Ask AI", func() {
		label.SetText("AI is thinking...")
	})

	content := container.NewVBox(
		label, button,
	)

	w.SetContent(content)

	w.ShowAndRun()
}
