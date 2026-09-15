package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("Dynamic Widgets")

	messages := container.NewVBox()

	addButton := widget.NewButton(
		"New Message",
		func() {
			messages.Add(
				widget.NewLabel("New Message"),
			)

			messages.Refresh()
		})

	content := container.NewVBox(
		addButton,
		messages,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
