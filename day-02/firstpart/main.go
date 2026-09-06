package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("AI Settings")
	w.Resize(fyne.NewSize(500, 400))
	w.SetFixedSize(true)

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter your name...")

	modelSelect := widget.NewSelect(
		[]string{
			"qwen2.5:3b",
			"gemma3:4b",
			"llama3.2",
		}, func(value string) {
			fmt.Println("Model:", value)
		},
	)

	modelSelect.SetSelected("qwen2.5:3b")

	tempratureLabel := widget.NewLabel("Temperature : 0.7")

	temprature := widget.NewSlider(0, 100)
	temprature.SetValue(0.7)

	temprature.OnChanged = func(f float64) {
		tempratureLabel.SetText(
			fmt.Sprintf("Temperature: %.2f", f),
		)
	}

	streming := widget.NewCheck(
		"Enable Streaming",
		func(b bool) {
			fmt.Println("Streaming:", b)
		},
	)

	responseType := widget.NewRadioGroup(
		[]string{
			"Text",
			"Json",
			"Markdown",
		}, func(s string) {
			fmt.Println("Response:", s)
		},
	)

	startButton := widget.NewButton("Start AI", func() {
		fmt.Println("========== AI SETTINGS ==========")
		fmt.Println("Name:", nameEntry.Text)
		fmt.Println("Model:", modelSelect.Selected)
		fmt.Println("Temperature:", temprature.Value)
		fmt.Println("Streaming:", streming.Checked)
		fmt.Println("=================================")
	})

	content := container.NewVBox(
		widget.NewLabel("AI Settings"),
		widget.NewLabel("Name:"),
		nameEntry,

		widget.NewLabel("Model:"),
		modelSelect,

		tempratureLabel,
		temprature,

		streming,

		widget.NewLabel("Response Type:"),
		responseType,

		startButton,
	)

	w.SetContent(content)

	w.ShowAndRun()
}
