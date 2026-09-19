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

	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			FOREIGN KEY (chat_id)
				REFERENCES chats(id)
				ON DELETE CASCADE	
		)
		`)

	if err != nil {
		return err
	}

	return nil
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

func (r *ChatRepository) CreateMessage(chat_id int, role string, content string) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO messages (chat_id, role, content) VALUES (?, ?, ?)",
		chat_id,
		role,
		content,
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

func (r *ChatRepository) GetMessages(chat_id int) ([]ChatMessage, error) {
	rows, err := r.db.Query(
		"SELECT id, chat_id, role, content FROM messages WHERE chat_id = ?",
		chat_id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []ChatMessage

	for rows.Next() {
		var message ChatMessage

		err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.Role,
			&message.Content,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
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

func (r *ChatRepository) DeletetMessages(chat_id int) error {
	_, err := r.db.Exec(
		"DELETE FROM messages WHERE chat_id = ?",
		chat_id,
	)

	return err
}
