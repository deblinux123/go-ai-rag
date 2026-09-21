package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	OllamaURL    string  `json:"ollama_url"`
	Model        string  `json:"model"`
	Temperature  float64 `json:"temperature"`
	Theme        string  `json:"theme"`
	SystemPrompt string  `json:"system_prompt"`
}

func DefaultConfig() Config {
	return Config{
		OllamaURL:    "http://localhost:11434",
		Model:        "qwen2.5:3b",
		Temperature:  0.7,
		Theme:        "dark",
		SystemPrompt: "You are a helpful AI assistant. Be concise, accurate, and clear.",
	}
}

type Store struct {
	mu   sync.RWMutex
	path string
	cfg  Config
}

func Load() (*Store, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("get user config directory: %w", err)
	}

	appDir := filepath.Join(dir, "GoAIRAG")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return nil, fmt.Errorf("create config directory: %w", err)
	}

	path := filepath.Join(appDir, "settings.json")
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("decode settings: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read settings: %w", err)
	}

	store := &Store{
		path: path,
		cfg:  cfg,
	}

	if err := store.Save(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) Get() Config {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.cfg
}

func (s *Store) Update(update func(*Config)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	update(&s.cfg)

	data, err := json.MarshalIndent(s.cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode settings: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("write settings: %w", err)
	}

	return nil
}

func (s *Store) Save() error {
	s.mu.RLock()
	data, err := json.MarshalIndent(s.cfg, "", "  ")
	s.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("encode settings: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("write settings: %w", err)
	}

	return nil
}

func (s *Store) DatabasePath() string {
	return filepath.Join(filepath.Dir(s.path), "conversations.db")
}
