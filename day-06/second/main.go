package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("Goroutine Test")

	status := widget.NewLabel("Ready")

	startButton := widget.NewButton("Start", func() {
		go func() {
			fmt.Println("Background task started.")

			time.Sleep(5 * time.Second)

			fmt.Println("Background task finished.")
		}()
	})

	content := container.NewVBox(
		status,
		startButton,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
