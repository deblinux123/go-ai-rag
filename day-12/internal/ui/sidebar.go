package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/deblinux123/go-ai-rag/day-12/internal/chat"
)

type Sidebar struct {
	app  *App
	list *widget.List

	chats []chat.Chat
}

func NewSidebar(app *App) fyne.CanvasObject {
	s := &Sidebar{
		app: app,
	}

	title := widget.NewLabelWithStyle(
		"Chats",
		fyne.TextAlignLeading,
		fyne.TextStyle{
			Bold: true,
		},
	)

	newChatButton := widget.NewButton(
		"+ New Chat",
		func() {
			s.createChat()
		},
	)

	s.list = widget.NewList(
		func() int {
			return len(s.chats)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Chat")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			if id >= len(s.chats) {
				return
			}

			label.SetText(
				fmt.Sprintf(
					"# %s",
					s.chats[id].Title,
				),
			)
		},
	)

	s.list.OnSelected = func(id widget.ListItemID) {
		if id >= len(s.chats) {
			return
		}

		s.app.selectedChatID = s.chats[id].ID
	}

	header := container.NewVBox(
		title,
		newChatButton,
		widget.NewSeparator(),
	)

	content := container.NewBorder(
		header,
		nil,
		nil,
		nil,
		s.list,
	)

	s.loadChats()

	return content
}

func (s *Sidebar) loadChats() {
	chats, err := s.app.chatService.GetChats()
	if err != nil {
		s.app.ShowError(err)
		return
	}

	s.chats = chats
	s.list.Refresh()
}

func (s *Sidebar) createChat() {
	id, err := s.app.chatService.CreateChat("New Chat")
	if err != nil {
		s.app.ShowError(err)
		return
	}

	s.app.selectedChatID = id

	s.loadChats()
}
