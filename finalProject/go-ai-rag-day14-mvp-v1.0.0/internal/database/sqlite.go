package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLite struct {
	DB *sql.DB
}

func NewSQLite(path string) (*SQLite, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	sqlite := &SQLite{DB: db}

	if err := sqlite.createTables(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create tables: %w", err)
	}

	return sqlite, nil
}

func (s *SQLite) createTables() error {
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS chats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (chat_id)
				REFERENCES chats(id)
				ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_messages_chat_id
			ON messages(chat_id);

		CREATE INDEX IF NOT EXISTS idx_chats_updated_at
			ON chats(updated_at);
	`)

	return err
}

func (s *SQLite) Close() error {
	return s.DB.Close()
}
