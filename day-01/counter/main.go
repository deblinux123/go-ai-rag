package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("com.example.com")

	w := a.NewWindow("Counter")

	count := 0

	label := widget.NewLabel("Count: 0")

	button := widget.NewButton("+", func() {
		count++

		label.SetText(fmt.Sprintf("Count: %d", count))
	})

	content := container.NewVBox(
		label, button,
	)

	w.SetContent(content)

	w.ShowAndRun()
}
