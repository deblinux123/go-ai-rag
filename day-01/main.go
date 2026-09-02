package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("Go AI RAG")

	w.SetContent(
		widget.NewLabel("Go AI RAG 🚀"),
	)

	w.ShowAndRun()
}
