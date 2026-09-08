package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AIConfig struct {
	Model        string
	SystemPrompt string
	Temperature  float64
	Streaming    bool
	ResponseType string
}

func main() {
	a := app.New()

	w := a.NewWindow("AI Configuration")

	title := canvas.NewText(
		"AI Configuration",
		color.RGBA{
			R: 50,
			G: 100,
			B: 200,
			A: 255,
		},
	)

	title.TextSize = 24

	title.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	modelLabel := widget.NewLabel("AI Model")
	modelSelect := widget.NewSelect(
		[]string{
			"qwen2.5:3b",
			"gemma3:4b",
			"llama3.2",
			"mistral",
		},
		nil,
	)

	modelSelect.SetSelected("qwen2.5:3b")

	promptLabel := widget.NewLabel("System Prompt")
	systemPrompt := widget.NewMultiLineEntry()

	systemPrompt.SetPlaceHolder("You are a helpfull AI asistant...")

	systemPrompt.SetMinRowsVisible(5)

	systemPrompt.SetText("You are a helpful AI assistant.")

	temperatureLabel := widget.NewLabel("Temperature: 0.7")
	temperatureSlider := widget.NewSlider(
		0.0,
		2.0,
	)

	temperatureSlider.Value = 0.7
	temperatureSlider.OnChangeEnded = func(f float64) {
		temperatureLabel.SetText(
			fmt.Sprintf("Temperature: %.1f", f),
		)
	}

	streamCheck := widget.NewCheck("Enable Streaming", nil)

	streamCheck.SetChecked(true)

	responseTypeLabel := widget.NewLabel("Response Type")

	responseType := widget.NewRadioGroup(
		[]string{
			"Text",
			"Json",
			"Markdown",
		},
		nil,
	)

	responseType.SetSelected("Text")

	statusLabel := widget.NewLabel("Status: Ready.")

	validationLabel := canvas.NewText(
		"",
		color.RGBA{
			R: 255,
			G: 0,
			B: 0,
			A: 255,
		},
	)

	progress := widget.NewProgressBar()

	progress.Hide()

	var saveButton *widget.Button

	validate := func() {
		if modelSelect.Selected == "" {
			saveButton.Disable()
			return
		}

		if systemPrompt.Text == "" {
			saveButton.Disable()
			return
		}

		if temperatureSlider.Value < 0 || temperatureSlider.Value > 2 {
			saveButton.Disable()
			return
		}

		saveButton.Enable()
	}

	saveConfig := func() {
		if modelSelect.Selected == "" {
			statusLabel.SetText("Plase select an AI Model.")
			return
		}

		if systemPrompt.Text == "" {
			statusLabel.SetText("System prompt is required.")
			return
		}

		config := AIConfig{
			Model:        modelSelect.Selected,
			SystemPrompt: systemPrompt.Text,
			Temperature:  temperatureSlider.Value,
			Streaming:    streamCheck.Checked,
			ResponseType: responseType.Selected,
		}

		progress.Show()
		progress.SetValue(0.0)

		statusLabel.SetText("Saving configuration...")

		progress.SetValue(0.5)

		statusLabel.SetText(
			fmt.Sprintf("Configuration saved: %s", config.Model),
		)

		validationLabel.Text = fmt.Sprintf(
			"Model: %s | Temperature: %.1f | Streaming: %t | Response: %s",
			config.Model,
			config.Temperature,
			config.Streaming,
			config.ResponseType,
		)

		validationLabel.Color = color.RGBA{
			R: 0,
			G: 160,
			B: 0,
			A: 255,
		}

		validationLabel.Refresh()

		progress.SetValue(0.1)

		progress.Hide()
	}

	saveButton = widget.NewButton("Save Configuration.", saveConfig)

	saveButton.Disable()

	modelSelect.OnChanged = func(s string) {
		validationLabel.Text = ""
		validationLabel.Refresh()

		statusLabel.SetText(
			"Model selected:" + s,
		)

		validate()
	}

	systemPrompt.OnChanged = func(s string) {
		if s == "" {
			validationLabel.Text = "System Prompt is required."

			validationLabel.Color = color.RGBA{
				R: 255,
				G: 0,
				B: 0,
				A: 255,
			}
		} else {
			validationLabel.Text = ""
			validationLabel.Refresh()
		}

		validationLabel.Refresh()

		validate()
	}

	streamCheck.OnChanged = func(b bool) {
		if b {
			statusLabel.SetText(
				"Streaming enabled.",
			)
		} else {
			statusLabel.SetText(
				"Streaming disabled.",
			)
		}
	}

	responseType.OnChanged = func(s string) {
		statusLabel.SetText(
			"Response Type:" + s,
		)
	}

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		modelLabel,
		modelSelect,

		promptLabel,
		systemPrompt,
		temperatureLabel,
		temperatureSlider,

		streamCheck,

		responseTypeLabel,
		responseType,

		widget.NewSeparator(),

		validationLabel,

		saveButton,
		progress,
		statusLabel,
	)

	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)
	w.SetContent(
		container.NewScroll(content),
	)
	w.ShowAndRun()
}
