package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"github.com/deblinux123/go-ai-rag/day-12/internal/ai"
	"github.com/deblinux123/go-ai-rag/day-12/internal/chat"
	"github.com/deblinux123/go-ai-rag/day-12/internal/database"
	"github.com/deblinux123/go-ai-rag/day-12/internal/ui"
)

func main() {
	db, err := database.NewSQLite("chat.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	chatService := chat.NewService(db.DB)

	ollamaClient := ai.NewClient(
		"http://localhost:11434",
	)

	fyneApp := app.New()

	application := ui.NewAPP(
		fyneApp,
		chatService,
		ollamaClient,
	)

	application.Run()
}
