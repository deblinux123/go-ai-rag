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

	w := a.NewWindow("Channel Demo")

	status := widget.NewLabel("Ready")

	progress := widget.NewProgressBar()
	progress.Max = 100

	var startButton *widget.Button

	startButton = widget.NewButton("Start", func() {
		status.SetText("Processing...")
		startButton.Disable()

		progressChan := make(chan int)

		go func() {
			for i := 0; i <= 100; i++ {
				time.Sleep(10 * time.Millisecond)
				progressChan <- i
			}
			close(progressChan)
		}()

		go func() {
			for value := range progressChan {
				fyne.Do(func() {
					progress.SetValue(float64(value))
				})
			}

			fyne.Do(func() {
				status.SetText("Done!")
				startButton.Enable()
			})
		}()
	})

	content := container.NewVBox(
		status,
		startButton,
		progress,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
