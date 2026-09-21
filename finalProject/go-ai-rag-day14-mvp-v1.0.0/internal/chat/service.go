package chat

import (
	"database/sql"
	"fmt"
	"strings"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateChat(title string) (int64, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		title = "New Chat"
	}

	result, err := s.db.Exec(
		"INSERT INTO chats (title) VALUES (?)",
		title,
	)
	if err != nil {
		return 0, fmt.Errorf("create chat: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get chat id: %w", err)
	}

	return id, nil
}

func (s *Service) UpdateChatTitle(id int64, title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil
	}

	_, err := s.db.Exec(
		"UPDATE chats SET title = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		title,
		id,
	)
	if err != nil {
		return fmt.Errorf("update chat title: %w", err)
	}

	return nil
}

func (s *Service) TouchChat(id int64) error {
	_, err := s.db.Exec(
		"UPDATE chats SET updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("touch chat: %w", err)
	}

	return nil
}

func (s *Service) GetChats() ([]Chat, error) {
	rows, err := s.db.Query(`
		SELECT id, title, created_at, updated_at
		FROM chats
		ORDER BY updated_at DESC, id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("get chats: %w", err)
	}
	defer rows.Close()

	var chats []Chat

	for rows.Next() {
		var c Chat
		if err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan chat: %w", err)
		}

		chats = append(chats, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate chats: %w", err)
	}

	return chats, nil
}

func (s *Service) GetChat(id int64) (*Chat, error) {
	var c Chat

	err := s.db.QueryRow(`
		SELECT id, title, created_at, updated_at
		FROM chats
		WHERE id = ?
	`, id).Scan(
		&c.ID,
		&c.Title,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get chat: %w", err)
	}

	return &c, nil
}

func (s *Service) DeleteChat(id int64) error {
	_, err := s.db.Exec(
		"DELETE FROM chats WHERE id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("delete chat: %w", err)
	}

	return nil
}

func (s *Service) CreateMessage(chatID int64, role, content string) (int64, error) {
	result, err := s.db.Exec(`
		INSERT INTO messages (chat_id, role, content)
		VALUES (?, ?, ?)
	`, chatID, role, content)
	if err != nil {
		return 0, fmt.Errorf("create message: %w", err)
	}

	if err := s.TouchChat(chatID); err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get message id: %w", err)
	}

	return id, nil
}

func (s *Service) GetMessages(chatID int64) ([]Message, error) {
	rows, err := s.db.Query(`
		SELECT id, chat_id, role, content, created_at
		FROM messages
		WHERE chat_id = ?
		ORDER BY id ASC
	`, chatID)
	if err != nil {
		return nil, fmt.Errorf("get messages: %w", err)
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
			&message.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate messages: %w", err)
	}

	return messages, nil
}
