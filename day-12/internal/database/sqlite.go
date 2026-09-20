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
		return nil, fmt.Errorf("Open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("Ping database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_key = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("Enable foreign keys: %w", err)
	}

	sqlite := &SQLite{
		DB: db,
	}

	if err := sqlite.createTable(); err != nil {
		db.Close()
		return nil, fmt.Errorf("Create tables: %w", err)
	}

	return sqlite, nil
}

func (s *SQLite) createTable() error {
	_, err := s.DB.Exec(`
	CREATE TABLE IF NOT EXISTS chats(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		chat_id INTEGER NOT NULL,
		role TEXT NOT NULL,
		content TEXT NO NULL,
		FOREIGN KEY (chat_id)
			REFERENCES chats(id)
			ON DELETE CASCADE
	);
	`)

	return err
}
func (s *SQLite) Close() error {
	return s.DB.Close()
}
