package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Model struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type OllamaResponse struct {
	Models []Model `json:"models"`
}

type Manager struct {
	app       fyne.App
	window    fyne.Window
	status    *widget.Label
	modelList *widget.List
	models    []Model
}

func formatSize(size int64) string {
	if size == 0 {
		return "0 B"
	}
	const unit = 1024
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func (m *Manager) loadModels() {
	m.status.SetText("Status: Connecting to Ollama...")

	go func() {
		resp, err := http.Get("http://localhost:11434/api/tags")
		if err != nil {
			fyne.Do(func() {
				m.status.SetText("Status: Connection failed. Is Ollama running?")
				dialog.ShowError(fmt.Errorf("failed to connect to Ollama: %w", err), m.window)
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fyne.Do(func() {
				m.status.SetText(fmt.Sprintf("Status: Ollama returned %s", resp.Status))
				dialog.ShowError(fmt.Errorf("unexpected status: %s", resp.Status), m.window)
			})
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fyne.Do(func() {
				m.status.SetText("Status: Failed to read response")
				dialog.ShowError(err, m.window)
			})
			return
		}

		var result OllamaResponse
		if err := json.Unmarshal(body, &result); err != nil {
			fyne.Do(func() {
				m.status.SetText("Status: Invalid JSON response")
				dialog.ShowError(err, m.window)
			})
			return
		}

		fyne.Do(func() {
			m.models = result.Models
			m.modelList.Refresh()
			m.status.SetText(fmt.Sprintf("Status: %d model(s) found", len(m.models)))
		})
	}()
}

func (m *Manager) pullModel(name string) {
	m.status.SetText(fmt.Sprintf("Status: Pulling %s... (this may take a while)", name))

	go func() {
		payload := fmt.Sprintf(`{"name": "%s", "stream": false}`, name)
		req, err := http.NewRequest("POST", "http://localhost:11434/api/pull", bytes.NewBufferString(payload))
		if err != nil {
			fyne.Do(func() { dialog.ShowError(err, m.window) })
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 0}
		resp, err := client.Do(req)
		if err != nil {
			fyne.Do(func() {
				m.status.SetText("Status: Pull failed")
				dialog.ShowError(err, m.window)
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fyne.Do(func() {
				m.status.SetText("Status: Pull failed")
				dialog.ShowError(fmt.Errorf("pull failed: %s", string(body)), m.window)
			})
			return
		}

		fyne.Do(func() {
			m.status.SetText(fmt.Sprintf("Status: Successfully pulled %s", name))
			m.loadModels()
		})
	}()
}

func (m *Manager) deleteModel(name string) {
	dialog.ShowConfirm("Delete Model", fmt.Sprintf("Are you sure you want to delete '%s'?", name), func(confirmed bool) {
		if !confirmed {
			return
		}

		m.status.SetText(fmt.Sprintf("Status: Deleting %s...", name))
		go func() {
			payload := fmt.Sprintf(`{"name": "%s"}`, name)
			req, err := http.NewRequest("DELETE", "http://localhost:11434/api/delete", bytes.NewBufferString(payload))
			if err != nil {
				fyne.Do(func() { dialog.ShowError(err, m.window) })
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fyne.Do(func() {
					m.status.SetText("Status: Delete failed")
					dialog.ShowError(err, m.window)
				})
				return
			}
			defer resp.Body.Close()

			fyne.Do(func() {
				if resp.StatusCode == http.StatusOK {
					m.status.SetText(fmt.Sprintf("Status: Deleted %s", name))
					m.loadModels()
				} else {
					m.status.SetText("Status: Delete failed")
					body, _ := io.ReadAll(resp.Body)
					dialog.ShowError(fmt.Errorf("delete failed: %s", string(body)), m.window)
				}
			})
		}()
	}, m.window)
}

// createListItem is the factory function for the list
func createListItem() fyne.CanvasObject {
	nameLabel := widget.NewLabel("Model Name")
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	sizeLabel := widget.NewLabel("Size")
	sizeLabel.Alignment = fyne.TextAlignTrailing

	deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {})
	deleteBtn.Importance = widget.DangerImportance

	// NewHBox guarantees the Objects slice order matches the arguments exactly
	return container.NewHBox(
		nameLabel,          // Index 0
		layout.NewSpacer(), // Index 1
		sizeLabel,          // Index 2
		deleteBtn,          // Index 3
	)
}

// updateListItem updates the factory-created item with actual data
func (m *Manager) updateListItem(id widget.ListItemID, item fyne.CanvasObject) {
	if id >= len(m.models) {
		return
	}

	model := m.models[id]
	hbox := item.(*fyne.Container)

	// Safely assert based on the known order from createListItem
	nameLabel := hbox.Objects[0].(*widget.Label)
	sizeLabel := hbox.Objects[2].(*widget.Label)
	deleteBtn := hbox.Objects[3].(*widget.Button)

	nameLabel.SetText(model.Name)
	sizeLabel.SetText(formatSize(model.Size))

	// Capture the model name for the closure
	currentModel := model.Name
	deleteBtn.OnTapped = func() {
		m.deleteModel(currentModel)
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Ollama Model Manager")
	w.Resize(fyne.NewSize(650, 500))

	manager := &Manager{
		app:    a,
		window: w,
	}

	manager.status = widget.NewLabel("Status: Ready")
	manager.status.TextStyle = fyne.TextStyle{Italic: true}

	refreshBtn := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), manager.loadModels)

	pullEntry := widget.NewEntry()
	pullEntry.SetPlaceHolder("e.g., llama3:8b")
	// Allow pressing "Enter" in the text field to trigger the pull
	pullEntry.OnSubmitted = func(s string) {
		if strings.TrimSpace(s) != "" {
			manager.pullModel(strings.TrimSpace(s))
			pullEntry.SetText("")
		}
	}

	pullBtn := widget.NewButtonWithIcon("Pull Model", theme.DownloadIcon(), func() {
		name := strings.TrimSpace(pullEntry.Text)
		if name == "" {
			dialog.ShowInformation("Input Required", "Please enter a model name to pull.", w)
			return
		}
		manager.pullModel(name)
		pullEntry.SetText("")
	})

	manager.modelList = widget.NewList(
		func() int {
			return len(manager.models)
		},
		createListItem,
		manager.updateListItem,
	)

	header := container.NewVBox(
		widget.NewLabelWithStyle("Ollama Model Manager", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		manager.status,
	)

	toolbar := container.NewHBox(
		refreshBtn,
		layout.NewSpacer(),
		pullEntry,
		pullBtn,
	)

	mainContent := container.NewBorder(
		container.NewVBox(header, container.NewPadded(toolbar)),
		nil,
		nil,
		nil,
		manager.modelList,
	)

	w.SetContent(mainContent)

	manager.loadModels()

	w.ShowAndRun()
}
