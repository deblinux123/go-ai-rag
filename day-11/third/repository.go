package main

import "database/sql"

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (r *ChatRepository) CreateTable() error {
	_, err := r.db.Exec(
		`CREATE TABLE IF NOT EXISTS chats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL
		)
	`)

	return err
}

func (r *ChatRepository) CreateChat(title string) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO chats (title) VALUES (?)",
		title,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ChatRepository) GetChats() ([]Chat, error) {
	rows, err := r.db.Query(
		"SELECT id, title FROM chats ORDER BY id",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chats []Chat
	for rows.Next() {

		var chat Chat

		err := rows.Scan(
			&chat.ID,
			&chat.Title,
		)

		if err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, err
}

func (r *ChatRepository) GetChat(id int) (*Chat, error) {
	var chat Chat

	err := r.db.QueryRow(
		"SELECT id, title FROM chats WHERE id = ?",
		id,
	).Scan(
		&chat.ID,
		&chat.Title,
	)

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepository) UpdateChat(id int, title string) error {
	_, err := r.db.Exec(
		"UPDATE chats SET title = ? WHERE id = ?",
		title,
		id,
	)

	return err
}

func (r *ChatRepository) DeleteChat(id int) error {
	_, err := r.db.Exec(
		"DELETE FROM chats WHERE id = ?",
		id,
	)

	return err
}
