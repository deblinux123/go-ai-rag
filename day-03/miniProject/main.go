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

	w := a.NewWindow("Fyne Dashboard")

	// header
	header := container.NewHBox(
		widget.NewLabel("🤖 Fyne Dashboard"),
		widget.NewButton("Settings", func() {}),
	)

	// sidebar
	sidebar := container.NewVBox(
		widget.NewButton("+ New Task", func() {}),
		widget.NewButton("Tasks", func() {}),
		widget.NewButton("Statistics", func() {}),
		widget.NewButton("Settings", func() {}),
	)

	taskLabe := widget.NewLabel("Tasks: 12")
	doneLabel := widget.NewLabel("Done: 8")
	pendingLabel := widget.NewLabel("Pending: 4")
	// stat
	stats := container.NewGridWithColumns(
		3,
		taskLabe,
		doneLabel,
		pendingLabel,
	)

	// loading
	loading := container.NewVBox(
		widget.NewLabel("Loading..."),
		widget.NewProgressBarInfinite(),
	)

	loading.Hide()

	// dashboard stack
	dashboardStack := container.NewStack(
		stats,
		loading,
	)

	// refresh button
	var refreshButton *widget.Button
	refreshButton = widget.NewButton("Refresh", func() {
		loading.Show()
		refreshButton.Disable()

		go func() {
			time.Sleep(3 * time.Second)
			taskLabe.SetText("Tasks: 15")
			doneLabel.SetText("Done: 11")
			pendingLabel.SetText("Pending: 4")

			loading.Hide()
			refreshButton.Enable()
		}()
	})

	// main content
	mainContent := container.NewCenter(

		container.NewVBox(
			widget.NewLabel("Dashboard"),
			dashboardStack,
			refreshButton,
		),
	)

	// footer
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

	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)
	w.SetContent(content)
	w.ShowAndRun()
}
