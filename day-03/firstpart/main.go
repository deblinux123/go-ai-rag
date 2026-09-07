package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("AI Desktop Assistant")
	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)

	header := container.NewHBox(
		widget.NewLabel("🤖 AI Desktop Assistant"),
		widget.NewButton("Settings", func() {}),
	)

	sidebar := container.NewVBox(
		widget.NewButton("+ New Chat", func() {}),
		widget.NewButton("Chat History", func() {}),
		widget.NewButton("Models", func() {}),
		widget.NewButton("Settings", func() {}),
	)

	quickActionLabel := widget.NewLabel("Quick Actions.")
	quickAction := container.NewGridWithColumns(
		3,
		widget.NewButton("New Chat", func() {}),
		widget.NewButton("Models", func() {}),
		widget.NewButton("Clear", func() {}),
	)

	response := widget.NewMultiLineEntry()
	response.SetText("AI response will appear here...")

	loading := container.NewVBox(
		widget.NewLabel("AI is thinking..."),
		widget.NewProgressBarInfinite(),
	)
	loading.Hide()

	responseStack := container.NewStack(
		response,
		loading,
	)
	askButton := widget.NewButton("Ask AI", func() {
		loading.Show()

		response.SetText("AI is thinking...")

		go func() {
			time.Sleep(3 * time.Second)

			loading.Hide()

			response.SetText(
				"Hello! 👋\n\nThis is a fake AI response.\n\nNext step: We will connect this application to Ollama.",
			)
		}()
	})

	mainContent := container.NewCenter(
		container.NewVBox(
			widget.NewLabel("Welcome to AI Desktop Assistant"),
			quickActionLabel,
			quickAction,
			responseStack,
			askButton,
		),
	)

	footer := container.NewHBox(
		widget.NewLabel("Status: Ready"),
	)

	content := container.NewBorder(
		header,
		footer,
		sidebar,
		nil,
		mainContent,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
