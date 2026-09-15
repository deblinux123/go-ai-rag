package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const ollamaURL = "http://localhost:11434"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

type Model struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type ModelsResponse struct {
	Models []Model `json:"models"`
}

type OllamaClient struct {
	BaseURL string
	Client  *http.Client
}

func NewOllamaClient() *OllamaClient {
	return &OllamaClient{
		BaseURL: ollamaURL,
		Client: &http.Client{
			Timeout: 3 * time.Minute,
		},
	}
}

func (c *OllamaClient) GetModels() ([]Model, error) {
	resp, err := c.Client.Get(c.BaseURL + "/api/tags")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama returned status: %s", resp.Status)
	}

	var result ModelsResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Models, nil
}

func (c *OllamaClient) Chat(model, prompt string) (string, error) {
	requestBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	data, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.BaseURL+"/api/chat",
		bytes.NewBuffer(data),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Content-type",
		"application/json",
	)

	resp, err := c.Client.Do(req)

	if err != nil {
		return "", err
	}

	var result ChatResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Message.Content, nil
}
