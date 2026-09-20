package chat

import (
	"database/sql"
	"fmt"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateChat(title string) (int64, error) {
	result, err := s.db.Exec(
		"INSERT INTO chat (title) VALUES (?)",
		title,
	)

	if err != nil {
		return 0, fmt.Errorf("Create Chat: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("Get Chat id: %w", err)
	}

	return id, nil
}

func (s *Service) GetChats() ([]Chat, error) {
	rows, err := s.db.Query(
		"SELECT id, title FROM chats ORDER BY id DESC",
	)

	if err != nil {
		return nil, fmt.Errorf("Get Chats: %w", err)
	}

	defer rows.Close()

	var chats []Chat

	for rows.Next() {
		var c Chat

		if err := rows.Scan(
			&c.ID,
			&c.Title,
		); err != nil {
			return nil, fmt.Errorf("Scan Chat: %w", err)
		}

		chats = append(chats, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Iterate Chats: %w", err)
	}

	return chats, nil
}

func (s *Service) GetChat(id int64) (*Chat, error) {
	var c Chat

	err := s.db.QueryRow(
		"SELECT id, title FROM chats WHERE id = ?",
		id,
	).Scan(
		&c.ID,
		&c.Title,
	)

	if err != nil {
		return nil, fmt.Errorf("Get Chat: %w", err)
	}

	return &c, nil
}

func (s *Service) DeleteChat(id int64) error {
	_, err := s.db.Exec(
		"DELETE FROM chats WHERE id = ?",
		id,
	)

	if err != nil {
		return fmt.Errorf("Delete Chat: %w", err)
	}

	return nil
}

func (s *Service) CreateMessage(chatID int64, role, content string) (int64, error) {
	result, err := s.db.Exec(
		"INSERT INTO messages (chat_id, role, content) VALUES (?, ?, ?)",
		chatID,
		role,
		content,
	)

	if err != nil {
		return 0, fmt.Errorf("Create Message: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("Get Message ID: %w", err)
	}

	return id, nil
}

func (s *Service) GetMessages(chatID int64) ([]Message, error) {
	rows, err := s.db.Query(
		"SELECT id, chat_id, role, content FROM messages, WHERE chat_id = ? ORDER BY id ASC",
		chatID,
	)

	if err != nil {
		return nil, fmt.Errorf("Get Messages: %w", err)
	}

	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var message Message

		if err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.Role,
			&message.Content,
		); err != nil {
			return nil, fmt.Errorf("Scan Message: %w", err)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Iterate Messages: %w", err)
	}

	return messages, nil
}
