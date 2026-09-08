package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	w := a.NewWindow("Login Demo")

	usernameLabel := widget.NewLabel("Username")
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Enter username ...")

	passwordLabel := widget.NewLabel("Password")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter password ...")

	statusLabel := widget.NewLabel("Status: Ready.")

	var submitButton *widget.Button
	submit := func() {
		if usernameEntry.Text == "" {
			statusLabel.SetText("Please enter username.")
			return
		}

		if passwordEntry.Text == "" {
			statusLabel.SetText("Please enter password.")
			return
		}

		statusLabel.SetText("Welcome " + usernameEntry.Text + " !")
		usernameEntry.SetText("")
		passwordEntry.SetText("")
	}

	validate := func() {
		if usernameEntry.Text == "" || len(passwordEntry.Text) < 4 {
			submitButton.Disable()
		} else {
			submitButton.Enable()
		}
	}

	submitButton = widget.NewButton("Submit", submit)
	submitButton.Disable()

	passwordEntry.OnSubmitted = func(s string) {
		submit()
	}

	previewLabel := widget.NewLabel("Hello " + usernameEntry.Text)
	passwordValidationLabel := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})

	usernameEntry.OnChanged = func(s string) {
		previewLabel.SetText("Hello, " + s)
		validate()
	}

	passwordEntry.OnChanged = func(s string) {
		if len(s) == 0 {
			passwordValidationLabel.Text = "Password is required."
			passwordValidationLabel.Color = color.RGBA{
				R: 255,
				G: 0,
				B: 0,
				A: 255,
			}
		} else if len(s) <= 3 {
			passwordValidationLabel.Text = "Password must be at least 4 characters."
			passwordValidationLabel.Color = color.RGBA{
				R: 255,
				G: 0,
				B: 0,
				A: 255,
			}
		} else {
			passwordValidationLabel.Text = "Password looks good!"
			passwordValidationLabel.Color = color.RGBA{
				R: 0,
				G: 180,
				B: 0,
				A: 255,
			}
		}

		passwordValidationLabel.Refresh()

		validate()
	}

	content := container.NewVBox(
		usernameLabel,
		usernameEntry,
		previewLabel,
		passwordLabel,
		passwordEntry,
		passwordValidationLabel,
		submitButton,

		statusLabel,
	)

	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(true)
	w.SetContent(content)
	w.ShowAndRun()
}
