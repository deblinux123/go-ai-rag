package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewSettingsView(app *App) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(
		"Settings",
		fyne.TextAlignLeading,
		fyne.TextStyle{
			Bold: true,
		},
	)

	modelSelect := widget.NewSelect(
		[]string{
			"qwen2.5:3b",
			"gemma3:4b",
			"llama3.2",
		},
		func(model string) {
			// Model selection will be connected
			// to Ollama in the next step.
		},
	)

	modelSelect.SetSelected("qwen2.5:3b")

	modelLabel := widget.NewLabel("Model")

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		modelLabel,
		modelSelect,
	)
}
