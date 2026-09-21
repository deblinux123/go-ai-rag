package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/deblinux123/go-ai-rag/day-14/internal/ai"
	"github.com/deblinux123/go-ai-rag/day-14/internal/settings"
)

type SettingsView struct {
	app    *App
	canvas fyne.CanvasObject

	urlEntry    *widget.Entry
	modelSelect *widget.Select
	tempSlider  *widget.Slider
	tempLabel   *widget.Label
	themeSelect *widget.Select
	prompt      *widget.Entry
	status      *widget.Label
}

func NewSettingsView(app *App) *SettingsView {
	cfg := app.Config()

	s := &SettingsView{app: app}

	s.urlEntry = widget.NewEntry()
	s.urlEntry.SetText(cfg.OllamaURL)
	s.urlEntry.SetPlaceHolder("http://localhost:11434")

	s.modelSelect = widget.NewSelect(
		[]string{cfg.Model},
		func(model string) {
			if model == "" {
				return
			}

			_ = app.settings.Update(func(c *settings.Config) {
				c.Model = model
			})
		},
	)
	s.modelSelect.SetSelected(cfg.Model)

	s.tempLabel = widget.NewLabel(fmt.Sprintf("%.1f", cfg.Temperature))
	s.tempSlider = widget.NewSlider(0, 2)
	s.tempSlider.Value = cfg.Temperature
	s.tempSlider.Step = 0.1
	s.tempSlider.OnChanged = func(value float64) {
		s.tempLabel.SetText(fmt.Sprintf("%.1f", value))

		_ = app.settings.Update(func(c *settings.Config) {
			c.Temperature = value
		})
	}

	s.themeSelect = widget.NewSelect(
		[]string{"dark", "light"},
		func(value string) {
			app.SetTheme(value)
		},
	)
	s.themeSelect.SetSelected(cfg.Theme)

	s.prompt = widget.NewMultiLineEntry()
	s.prompt.SetText(cfg.SystemPrompt)
	s.prompt.SetMinRowsVisible(5)
	s.prompt.Wrapping = fyne.TextWrapWord

	s.status = widget.NewLabel("Settings are saved automatically.")

	s.urlEntry.OnChanged = func(value string) {
		value = strings.TrimRight(strings.TrimSpace(value), "/")

		_ = app.settings.Update(func(c *settings.Config) {
			c.OllamaURL = value
		})

		app.UpdateAIClientURL("http://localhost:11434")
	}

	s.prompt.OnChanged = func(value string) {
		_ = app.settings.Update(func(c *settings.Config) {
			c.SystemPrompt = value
		})
	}

	testButton := widget.NewButtonWithIcon(
		"Test Ollama",
		theme.SearchIcon(),
		func() {
			s.status.SetText("Checking Ollama…")

			go func() {
				app.aiClient.SetBaseURL(app.Config().OllamaURL)
				models, err := app.aiClient.GetModels()

				fyne.Do(func() {
					if err != nil {
						s.status.SetText("Connection failed")
						app.ShowError(err)
						return
					}

					s.SetModels(models)
					s.status.SetText(fmt.Sprintf("%d model(s) available", len(models)))
				})
			}()
		},
	)

	title := widget.NewLabelWithStyle(
		"Settings",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	form := widget.NewForm(
		widget.NewFormItem("Ollama URL", s.urlEntry),
		widget.NewFormItem("Model", s.modelSelect),
		widget.NewFormItem(
			"Temperature",
			container.NewBorder(
				nil,
				nil,
				nil,
				s.tempLabel,
				s.tempSlider,
			),
		),
		widget.NewFormItem("Theme", s.themeSelect),
		widget.NewFormItem("System prompt", s.prompt),
	)

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		form,
		testButton,
		s.status,
	)

	shell := container.NewVScroll(container.NewPadded(content))
	return s.withCanvas(shell)
}

func (s *SettingsView) withCanvas(canvas fyne.CanvasObject) *SettingsView {
	s.canvas = canvas
	return s
}

func (s *SettingsView) CanvasObject() fyne.CanvasObject {
	return s.canvas
}

func (s *SettingsView) SetModels(models []ai.Model) {
	if len(models) == 0 {
		return
	}

	options := make([]string, 0, len(models))
	selected := s.app.Config().Model

	for _, model := range models {
		options = append(options, model.Name)
	}

	s.modelSelect.SetOptions(options)

	found := false
	for _, option := range options {
		if option == selected {
			found = true
			break
		}
	}

	if found {
		s.modelSelect.SetSelected(selected)
		return
	}

	s.modelSelect.SetSelected(options[0])

	_ = s.app.settings.Update(func(c *settings.Config) {
		c.Model = options[0]
	})
}
