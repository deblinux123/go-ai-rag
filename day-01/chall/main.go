package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// creating the app
	a := app.NewWithID("com.example.helloai")

	// creating the window
	w := a.NewWindow("Hello AI")

	// create the label
	label := widget.NewLabel("Hello AI")

	// top button
	top := container.NewHBox(
		widget.NewLabel("AI"),
		widget.NewButton("Settings", func() {}),
	)

	responseLabel := widget.NewLabel("Response.")

	// create the response
	response := container.NewVBox(
		responseLabel,
		widget.NewButton("Ask", func() {
			responseLabel.SetText("Helo how can i help you.")
		}),
	)

	continent := container.NewVBox(
		label, top, response,
	)

	w.SetContent(continent)

	w.ShowAndRun()

}
