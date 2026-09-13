package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Model struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type OllamaResponse struct {
	Models []Model `json:"models"`
}

func main() {
	a := app.New()

	w := a.NewWindow("Ollama Model Manager")

	status := widget.NewLabel("Status: Ready")

	modelList := widget.NewList(
		func() int {
			return 0
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Model")
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {

		},
	)

	var getModelsButton *widget.Button
	getModelsButton = widget.NewButton("Get Ollama Models", func() {
		getModelsButton.OnTapped = func() {
			getModelsButton.Disable()
			status.SetText("Status: Connecting to Ollama...")

			go func() {
				resp, err := http.Get("http://localhost:11434/api/tags")

				if err != nil {
					fyne.Do(func() {
						status.SetText("Status: Connection failed.")
						getModelsButton.Enable()
					})
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					fyne.Do(func() {
						status.SetText(
							fmt.Sprintf("Status: Ollama returned %s", resp.Status),
						)
						getModelsButton.Enable()
					})

					return
				}

				body, err := io.ReadAll(resp.Body)

				if err != nil {
					fyne.Do(func() {
						status.SetText("Status: Failed to read response")
						getModelsButton.Enable()
					})
					return
				}

				var result OllamaResponse

				err = json.Unmarshal(body, &result)

				if err != nil {
					fyne.Do(func() {
						status.SetText("Status: Invalid JSON response")
						getModelsButton.Enable()
					})
					return
				}

				models := result.Models

				fyne.Do(func() {
					modelList.Length = func() int {
						return len(models)
					}

					modelList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
						label := item.(*widget.Label)
						if id < len(models) {
							label.SetText(models[id].Name)
						}
					}

					modelList.Refresh()

					status.SetText(
						fmt.Sprintf("Status: %d model(s) found", len(models)),
					)
					getModelsButton.Enable()
				})
			}()
		}
	})

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Ollama Models"),
			status,
			getModelsButton,
		),
		nil,
		nil,
		nil,
		modelList,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
