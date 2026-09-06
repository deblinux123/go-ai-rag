package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AIConfig struct {
	Name         string
	SystemPropmt string
	Model        string
	Temerature   float64
	Streaming    bool
	ResponseType string
}

func main() {
	a := app.New()

	w := a.NewWindow("AI Settings")
	w.Resize(fyne.NewSize(500, 500))
	w.SetFixedSize(true)

	// -------------------------
	// Name
	// -------------------------

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter your name...")

	// -------------------------
	// MultiLineEntry
	// -------------------------

	systemPrompLabel := widget.NewLabel("System Prompt:")

	multiLineEntry := widget.NewMultiLineEntry()
	multiLineEntry.SetPlaceHolder("Enter your system prompt.")

	// -------------------------
	// Model
	// -------------------------

	modelSelect := widget.NewSelect(
		[]string{
			"qwen2.5:3b",
			"gemma3:4b",
			"llama3.2",
		},
		func(value string) {
			fmt.Println("Model:", value)
		},
	)

	modelSelect.SetSelected("qwen2.5:3b")

	// -------------------------
	// Temperature
	// -------------------------

	temperatureLabel := widget.NewLabel("Temperature: 0.70")

	temperature := widget.NewSlider(0, 2)

	temperature.SetValue(0.7)

	temperature.OnChanged = func(value float64) {
		temperatureLabel.SetText(
			fmt.Sprintf("Temperature: %.2f", value),
		)
	}

	// -------------------------
	// Streaming
	// -------------------------

	streaming := widget.NewCheck(
		"Enable Streaming",
		func(value bool) {
			fmt.Println("Streaming:", value)
		},
	)

	// -------------------------
	// Response Type
	// -------------------------

	responseType := widget.NewRadioGroup(
		[]string{
			"Text",
			"JSON",
			"Markdown",
		},
		func(value string) {
			fmt.Println("Response:", value)
		},
	)

	// -------------------------
	// ProgressBar
	// -------------------------

	progress := widget.NewProgressBar()

	// -------------------------
	// Start AI
	// -------------------------

	multiLineResult := widget.NewMultiLineEntry()
	multiLineResult.Disable()

	startButton := widget.NewButton("Start AI", func() {

		config := AIConfig{
			Name:         nameEntry.Text,
			SystemPropmt: multiLineEntry.Text,
			Model:        modelSelect.Selected,
			Temerature:   temperature.Value,
			Streaming:    streaming.Checked,
			ResponseType: responseType.Selected,
		}

		multiLineResult.Enable()

		multiLineResult.SetText(
			fmt.Sprintf(
				"Name: %s\n\n"+
					"System Prompt: %s\n\n"+
					"Model: %s\n"+
					"Temperature: %.2f\n"+
					"Streaming: %t\n"+
					"Response Type: %s",
				config.Name,
				config.SystemPropmt,
				config.Model,
				config.Temerature,
				config.Streaming,
				config.ResponseType,
			),
		)

		go func() {
			for i := 0.0; i <= 1.0; i += 0.01 {
				time.Sleep(25 * time.Millisecond)
				progress.SetValue(i)
			}
		}()
	})

	// -------------------------
	// Layout
	// -------------------------

	content := container.NewVBox(

		widget.NewLabel("AI Settings"),

		widget.NewLabel("Name:"),
		nameEntry,

		systemPrompLabel,
		multiLineEntry,

		widget.NewLabel("Model:"),
		modelSelect,

		temperatureLabel,
		temperature,

		streaming,

		widget.NewLabel("Response Type:"),
		responseType,

		startButton,

		progress,

		multiLineResult,
	)

	w.SetContent(content)

	w.ShowAndRun()
}
