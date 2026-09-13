package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Post struct {
	UserId int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	a := app.New()

	w := a.NewWindow("HTTP GET Practice")

	urlEntry := widget.NewEntry()
	urlEntry.SetText("https://jsonplaceholder.typicode.com/posts/1")

	status := widget.NewLabel("Status: Ready")

	response := widget.NewMultiLineEntry()
	response.Wrapping = fyne.TextWrapWord
	response.Disable()

	getButton := widget.NewButton("GET", nil)

	getButton.OnTapped = func() {
		getButton.Disable()
		status.SetText("Status: Sending request...")
		response.SetText("")

		url := urlEntry.Text

		go func() {
			resp, err := http.Get(url)

			if err != nil {
				fyne.Do(func() {
					status.SetText("Status: Request failed")
					response.SetText(err.Error())
					getButton.Enable()
				})

				return
			}

			defer resp.Body.Close()

			var post Post
			err = json.NewDecoder(resp.Body).Decode(&post)

			if err != nil {
				fyne.Do(func() {
					status.SetText("Status: Failed to read response")
					response.SetText(err.Error())
					getButton.Enable()
				})
				return
			}

			result := fmt.Sprintf(
				"User ID: %d\n\nID: %d\n\nTitle: %s\n\nBody: %s\n\n",
				post.UserId,
				post.ID,
				post.Title,
				post.Body,
			)

			fyne.Do(func() {
				status.SetText(
					fmt.Sprintf("Status: %s", resp.Status),
				)

				response.SetText(result)
				getButton.Enable()
			})
		}()
	}

	content := container.NewVBox(
		widget.NewLabel("HTTP GET Practice"),

		widget.NewLabel("URL:"),
		urlEntry,
		getButton,
		status,
		widget.NewLabel("Response:"),
		response,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
