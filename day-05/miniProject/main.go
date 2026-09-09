package main

import (
	"fmt"
	"time"

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
	w := a.NewWindow("AI State Dashboard")

	// --------------------------------------------------
	// Application State
	// --------------------------------------------------

	state := AppState{
		Model:       "qwen2.5:3b",
		Temperature: 0.7,
		Streaming:   true,
	}

	// --------------------------------------------------
	// Model Binding
	// --------------------------------------------------

	modelBinding := binding.NewString()
	modelBinding.Set(state.Model)

	modelLabel := widget.NewLabelWithData(modelBinding)

	modelSelect := widget.NewSelectWithData(
		[]string{
			"qwen2.5:3b",
			"gemma3:4b",
			"llama3.2",
		},
		modelBinding,
	)

	// --------------------------------------------------
	// Temperature Binding
	// --------------------------------------------------

	temperatureBinding := binding.NewFloat()
	temperatureBinding.Set(state.Temperature)

	temperatureLabelBinding := binding.NewString()
	temperatureLabelBinding.Set(
		fmt.Sprintf("%.1f", state.Temperature),
	)

	temperatureLabel := widget.NewLabelWithData(
		temperatureLabelBinding,
	)

	temperatureSlider := widget.NewSliderWithData(
		0.0,
		2.0,
		temperatureBinding,
	)

	// --------------------------------------------------
	// Streaming Binding
	// --------------------------------------------------

	streamingBinding := binding.NewBool()
	streamingBinding.Set(state.Streaming)

	streamingCheck := widget.NewCheckWithData(
		"Streaming",
		streamingBinding,
	)

	streamingLabelBinding := binding.NewString()
	streamingLabelBinding.Set(
		fmt.Sprintf("%t", state.Streaming),
	)

	streamingLabel := widget.NewLabelWithData(
		streamingLabelBinding,
	)

	// --------------------------------------------------
	// Status Binding
	// --------------------------------------------------

	statusBinding := binding.NewString()
	statusBinding.Set("✓ Synced")

	statusLabel := widget.NewLabelWithData(statusBinding)

	// --------------------------------------------------
	// Animation
	// --------------------------------------------------

	activity := widget.NewActivity()
	activity.Hide()

	// Show animation while state is updating.
	setUpdating := func() {
		statusBinding.Set("◌ Updating...")

		activity.Show()
		activity.Start()

		// Stop animation after a short delay.
		time.AfterFunc(700*time.Millisecond, func() {
			fyne.Do(func() {
				activity.Stop()
				activity.Hide()
				statusBinding.Set("✓ Synced")
			})
		})
	}

	// --------------------------------------------------
	// Model Observer
	// --------------------------------------------------

	modelBinding.AddListener(
		binding.NewDataListener(func() {
			value, _ := modelBinding.Get()

			state.Model = value

			setUpdating()

			fmt.Println(
				"Model:",
				state.Model,
			)
		}),
	)

	// --------------------------------------------------
	// Temperature Observer
	// --------------------------------------------------

	temperatureBinding.AddListener(
		binding.NewDataListener(func() {
			value, _ := temperatureBinding.Get()

			state.Temperature = value

			temperatureLabelBinding.Set(
				fmt.Sprintf("%.1f", value),
			)

			setUpdating()

			fmt.Println(
				"Temperature:",
				state.Temperature,
			)
		}),
	)

	// --------------------------------------------------
	// Streaming Observer
	// --------------------------------------------------

	streamingBinding.AddListener(
		binding.NewDataListener(func() {
			value, _ := streamingBinding.Get()

			state.Streaming = value

			streamingLabelBinding.Set(
				fmt.Sprintf("%t", value),
			)

			setUpdating()

			fmt.Println(
				"Streaming:",
				state.Streaming,
			)
		}),
	)

	// --------------------------------------------------
	// Print State Button
	// --------------------------------------------------

	printStateButton := widget.NewButton(
		"Print State",
		func() {
			fmt.Printf(
				"Model=%s | Temperature=%.1f | Streaming=%t\n",
				state.Model,
				state.Temperature,
				state.Streaming,
			)
		},
	)

	// --------------------------------------------------
	// Status Section
	// --------------------------------------------------

	statusSection := container.NewHBox(
		activity,
		statusLabel,
	)

	// --------------------------------------------------
	// AI State Dashboard
	// --------------------------------------------------

	dashboard := container.NewVBox(
		widget.NewLabel("🤖 AI State Dashboard"),

		widget.NewSeparator(),

		widget.NewLabel("Model"),
		modelSelect,
		modelLabel,

		widget.NewLabel("Temperature"),
		temperatureSlider,
		temperatureLabel,

		widget.NewLabel("Streaming"),
		streamingCheck,
		streamingLabel,

		widget.NewSeparator(),

		widget.NewLabel("AI STATE"),
		statusSection,

		widget.NewSeparator(),

		printStateButton,
	)

	// --------------------------------------------------
	// Window
	// --------------------------------------------------

	w.Resize(fyne.NewSize(600, 500))
	w.SetFixedSize(true)
	w.SetContent(dashboard)
	w.ShowAndRun()
}
