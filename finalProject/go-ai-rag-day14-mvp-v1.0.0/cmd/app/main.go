package main

import (
	"log"

	"fyne.io/fyne/v2/app"

	"github.com/deblinux123/go-ai-rag/day-14/internal/ai"
	"github.com/deblinux123/go-ai-rag/day-14/internal/chat"
	"github.com/deblinux123/go-ai-rag/day-14/internal/database"
	"github.com/deblinux123/go-ai-rag/day-14/internal/settings"
	"github.com/deblinux123/go-ai-rag/day-14/internal/ui"
)

func main() {
	config, err := settings.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewSQLite(config.DatabasePath())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	chatService := chat.NewService(db.DB)
	ollamaClient := ai.NewClient(config.Get().OllamaURL)

	fyneApp := app.NewWithID("com.deblinux123.goairag")

	application := ui.NewApp(
		fyneApp,
		chatService,
		ollamaClient,
		config,
	)

	application.Run()
}
