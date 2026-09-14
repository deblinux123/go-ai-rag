package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const ollamaURL = "http://localhost:11434"

type Model struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type ModelsResponse struct {
	Models []Model `json:"models"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

type OllamaClinet struct {
	BaseURL string
	Client  *http.Client
}

func NewOllamaClient() *OllamaClinet {
	return &OllamaClinet{
		BaseURL: ollamaURL,
		Client: &http.Client{
			Timeout: 3 * time.Minute,
		},
	}
}

func (c *OllamaClinet) GetModels() ([]Model, error) {
	resp, err := c.Client.Get(c.BaseURL + "/api/tags")

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Ollama: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"Ollama returned %s: %s",
			resp.Status,
			strings.TrimSpace(string(body)),
		)
	}

	var result ModelsResponse

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, fmt.Errorf(
			"Failed to decode models response: %w",
			err,
		)
	}

	return result.Models, nil
}

func (c *OllamaClinet) Chat(model, prompt string) (string, error) {
	requestBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonBody, err := json.Marshal(requestBody)

	if err != nil {
		return "", fmt.Errorf(
			"Failed to encode request: %w",
			err,
		)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.BaseURL+"/api/chat",
		bytes.NewBuffer(jsonBody),
	)

	if err != nil {
		return "", fmt.Errorf(
			"Failed to create request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := c.Client.Do(req)

	if err != nil {
		return "", fmt.Errorf(
			"Failed to connect ollama: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return "", fmt.Errorf(
			"Ollama retured %s: %s",
			resp.Status,
			strings.TrimSpace(string(body)),
		)
	}

	var result ChatResponse

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return "", fmt.Errorf(
			"Failed to decde chat response: %w",
			err,
		)
	}

	return result.Message.Content, nil
}

type App struct {
	client *OllamaClinet
	window fyne.Window

	modelSelect *widget.Select
	prompt      *widget.Entry
	response    *widget.Entry

	status *widget.Label

	refreshButton *widget.Button
	sendButton    *widget.Button
}

func NewApp() *App {
	return &App{
		client: NewOllamaClient(),
	}
}

func (a *App) loadModels() {
	a.status.SetText("Status: Loading Models...")
	a.refreshButton.Disable()

	go func() {
		models, err := a.client.GetModels()

		if err != nil {
			fyne.Do(func() {
				a.status.SetText("Status: Failed to load models")

				dialog.ShowError(
					err,
					a.window,
				)

				a.refreshButton.Enable()
			})

			return
		}

		names := make([]string, 0, len(models))

		for _, model := range models {
			names = append(names, model.Name)
		}

		fyne.Do(func() {
			if len(names) == 0 {
				a.modelSelect.Options = nil
				a.modelSelect.SetSelected("")
				a.status.SetText("Status: No Model found.")
			} else {
				a.modelSelect.Options = names
				a.modelSelect.SetSelected(names[0])

				a.status.SetText(
					fmt.Sprintf(
						"Status: %d Model(s) found",
						len(names),
					),
				)
			}

			a.refreshButton.Enable()
		})
	}()
}

func (a *App) sendChat() {
	model := strings.TrimSpace(a.modelSelect.Selected)
	prompt := strings.TrimSpace(a.prompt.Text)

	if model == "" {
		dialog.ShowInformation(
			"Mode Required",
			"Please select an Ollama model.",
			a.window,
		)

		return
	}

	if prompt == "" {
		dialog.ShowInformation(
			"Prompt Required",
			"Please enter a prompt",
			a.window,
		)

		return
	}

	a.sendButton.Disable()
	a.refreshButton.Disable()

	a.status.SetText(
		fmt.Sprintf(
			"Status: %s is thinking...",
			model,
		),
	)

	a.response.SetText("")

	go func() {
		result, err := a.client.Chat(
			model,
			prompt,
		)

		if err != nil {
			fyne.Do(func() {
				a.status.SetText("Status: Chat Failed")

				dialog.ShowError(
					err,
					a.window,
				)

				a.sendButton.Enable()
				a.refreshButton.Enable()
			})

			return
		}

		fyne.Do(func() {
			a.response.SetText(result)

			a.status.SetText(
				fmt.Sprintf(
					"Status: Response recived from %s",
					model,
				),
			)

			a.sendButton.Enable()
			a.refreshButton.Enable()
		})
	}()
}

func (a *App) buildUI() {
	title := widget.NewLabelWithStyle(
		"Ollama Chat",
		fyne.TextAlignLeading,
		fyne.TextStyle{
			Bold: true,
		},
	)

	a.status = widget.NewLabel("Status: Ready")

	modelLabel := widget.NewLabel("Model")

	a.modelSelect = widget.NewSelect(
		nil,
		nil,
	)

	a.modelSelect.PlaceHolder = "Select Ollama model"

	a.refreshButton = widget.NewButtonWithIcon(
		"Refresh",
		theme.ViewRefreshIcon(),
		a.loadModels,
	)

	modelRow := container.NewBorder(
		nil,
		nil,
		modelLabel,
		a.refreshButton,
		a.modelSelect,
	)

	promptLabel := widget.NewLabel("Prompt")

	a.prompt = widget.NewMultiLineEntry()
	a.prompt.SetPlaceHolder(
		"Ask somthing...",
	)

	a.prompt.Wrapping = fyne.TextWrapWord

	a.sendButton = widget.NewButtonWithIcon(
		"Send",
		theme.MailSendIcon(),
		a.sendChat,
	)

	responseLabel := widget.NewLabel("Response")
	a.response = widget.NewMultiLineEntry()
	a.response.Wrapping = fyne.TextWrapWord
	a.response.Disable()

	content := container.NewVBox(
		title,
		a.status,
		widget.NewSeparator(),

		modelRow,
		promptLabel,
		a.prompt,
		a.sendButton,
		responseLabel,
		a.response,
	)

	a.window.SetContent(content)
}

func main() {
	fyneApp := app.New()

	window := fyneApp.NewWindow("Ollama")

	window.Resize(fyne.NewSize(800, 600))

	applicaiton := NewApp()

	applicaiton.window = window

	applicaiton.buildUI()

	applicaiton.loadModels()

	window.ShowAndRun()
}
