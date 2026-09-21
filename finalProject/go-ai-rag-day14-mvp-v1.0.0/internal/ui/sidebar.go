package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/deblinux123/go-ai-rag/day-14/internal/chat"
)

type Sidebar struct {
	app  *App
	list *widget.List

	chats []chat.Chat
}

func NewSidebar(app *App) *Sidebar {
	s := &Sidebar{app: app}

	s.list = widget.NewList(
		func() int {
			return len(s.chats)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Chat")
			label.Wrapping = fyne.TextWrapWord

			deleteButton := widget.NewButtonWithIcon(
				"",
				theme.DeleteIcon(),
				nil,
			)
			deleteButton.Importance = widget.LowImportance

			return container.NewBorder(
				nil,
				nil,
				nil,
				deleteButton,
				label,
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			row := obj.(*fyne.Container)
			label := row.Objects[0].(*widget.Label)
			deleteButton := row.Objects[len(row.Objects)-1].(*widget.Button)

			if int(id) >= len(s.chats) {
				return
			}

			currentID := s.chats[id].ID
			label.SetText(fmt.Sprintf("# %s", s.chats[id].Title))

			deleteButton.OnTapped = func() {
				dialog.ShowConfirm(
					"Delete conversation",
					"Delete this conversation and all of its messages?",
					func(ok bool) {
						if !ok {
							return
						}

						if err := s.app.chatService.DeleteChat(currentID); err != nil {
							s.app.ShowError(err)
							return
						}

						if s.app.selectedChatID == currentID {
							s.app.selectedChatID = 0
							s.app.chatView.Clear()
						}

						s.Reload()
					},
					s.app.window,
				)
			}
		},
	)

	s.list.HideSeparators = true
	s.list.OnSelected = func(id widget.ListItemID) {
		if int(id) >= len(s.chats) {
			return
		}

		s.app.SelectChat(s.chats[id].ID)
	}

	s.Reload()

	return s
}

func (s *Sidebar) CanvasObject() fyne.CanvasObject {
	title := widget.NewLabelWithStyle(
		"Conversations",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	newChat := widget.NewButtonWithIcon(
		"New chat",
		theme.ContentAddIcon(),
		func() {
			s.app.CreateChat()
		},
	)
	newChat.Importance = widget.HighImportance

	deleteChat := widget.NewButtonWithIcon(
		"Delete",
		theme.DeleteIcon(),
		func() {
			s.app.DeleteSelectedChat()
		},
	)
	deleteChat.Importance = widget.LowImportance

	header := container.NewBorder(
		nil,
		nil,
		nil,
		deleteChat,
		title,
	)

	return container.NewBorder(
		container.NewVBox(
			header,
			newChat,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		s.list,
	)
}

func (s *Sidebar) Reload() {
	chats, err := s.app.chatService.GetChats()
	if err != nil {
		s.app.ShowError(err)
		return
	}

	s.chats = chats
	s.list.Refresh()
}

func (s *Sidebar) SelectByID(id int64) {
	for i, c := range s.chats {
		if c.ID == id {
			s.list.Select(widget.ListItemID(i))
			return
		}
	}
}
