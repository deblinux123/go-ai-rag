package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type AppState struct {
	Model       string
	Temperature float64
	Streaming   bool
}

func main() {
	a := app.New()

	w := a.NewWindow("AI State")

	state := AppState{
		Model:       "qwen2.5:3b",
		Temperature: 0.7,
		Streaming:   true,
	}

	// String binding
	modelBinding := binding.NewString()
	modelBinding.Set(state.Model)

	// binding the float type
	temperatureBinding := binding.NewFloat()
	temperatureBinding.Set(state.Temperature)

	temperatureLabelBinding := binding.NewString()

	temperatureLabelBinding.Set(
		fmt.Sprintf("%.1f", state.Temperature),
	)

	temperatureLabel := widget.NewLabelWithData(
		temperatureLabelBinding,
	)

	temperatureBinding.AddListener(binding.NewDataListener(func() {
		value, _ := temperatureBinding.Get()

		temperatureLabelBinding.Set(
			fmt.Sprintf("%.1f", value),
		)
	}))
	// slider
	temperatureSlider := widget.NewSliderWithData(
		0.0,
		2.0,
		temperatureBinding,
	)

	// Label connected to binding
	modelLabel := widget.NewLabelWithData(modelBinding)

	// Select connected to the same binding
	modelSelect := widget.NewSelectWithData(
		[]string{
			"qwen2.5:3b",
			"gemma3:4b",
			"llama3.2",
		},
		modelBinding,
	)

	streamingBinding := binding.NewBool()
	streamingBinding.Set(state.Streaming)

	streamingChecked := widget.NewCheckWithData(
		"Streaming",
		streamingBinding,
	)

	streamingLabelBinding := binding.NewString()
	streamingLabelBinding.Set(
		fmt.Sprintf("%t", state.Streaming),
	)

	streamingBinding.AddListener(binding.NewDataListener(func() {
		value, _ := streamingBinding.Get()

		streamingLabelBinding.Set(
			fmt.Sprintf("%t", value),
		)
	}))

	streamingLabel := widget.NewLabelWithData(
		streamingLabelBinding,
	)

	content := container.NewVBox(
		widget.NewLabel("Model:"),
		modelLabel,
		modelSelect,
		widget.NewLabel("Temperature:"),
		temperatureLabel,
		temperatureSlider,
		widget.NewLabel("Streaming:"),
		streamingChecked,
		streamingLabel,
	)

	w.Resize(fyne.NewSize(600, 400))
	w.SetFixedSize(true)
	w.SetContent(content)
	w.ShowAndRun()
}
