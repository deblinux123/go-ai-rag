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

	w := a.NewWindow("AI Assistant")

	status := widget.NewLabel("Ready")

	response := widget.NewMultiLineEntry()
	response.SetPlaceHolder("AI response will appear here...")
	response.Wrapping = fyne.TextWrapWord
	response.Disable()

	var startButton *widget.Button
	startButton = widget.NewButton("Start", func() {
		startButton.OnTapped = func() {
			startButton.Disable()
			status.SetText("Thinking...")

			response.SetText("")

			responseChan := make(chan string)

			go func() {
				words := []string{
					"Hello! ",
					"I am ",
					"your ",
					"AI ",
					"assistant. ",
					"I can ",
					"help you ",
					"with your ",
					"questions.",
				}

				for _, word := range words {
					time.Sleep(300 * time.Millisecond)
					responseChan <- word
				}

				close(responseChan)
			}()

			go func() {
				fullResponse := ""

				for chunk := range responseChan {
					fullResponse += chunk

					fyne.Do(func() {
						response.SetText(fullResponse)
					})
				}

				fyne.Do(func() {
					status.SetText("Done!")
					startButton.Enable()
				})
			}()
		}
	})

	content := container.NewVBox(
		widget.NewLabel("AI Assistant"),
		status,
		startButton,
		response,
	)

	w.Resize(fyne.NewSize(600, 400))
	w.SetFixedSize(true)
	w.SetContent(content)
	w.ShowAndRun()
}
