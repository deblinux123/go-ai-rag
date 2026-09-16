package ollama

import (
	"bufio"
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
			Timeout: 10 * time.Minute,
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
		return nil, fmt.Errorf("ollama returned status: %s", resp.Status)
	}

	var result ModelsResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Models, nil
}

func (c *OllamaClient) StreamChat(
	model string,
	prompt string,
) (<-chan string, <-chan error) {

	tokens := make(chan string)
	errors := make(chan error, 1)

	go func() {
		defer close(tokens)
		defer close(errors)

		requestBody := ChatRequest{
			Model: model,
			Messages: []Message{
				{
					Role:    "user",
					Content: prompt,
				},
			},
			Stream: true,
		}

		data, err := json.Marshal(requestBody)
		if err != nil {
			errors <- err
			return
		}

		req, err := http.NewRequest(
			http.MethodPost,
			c.BaseURL+"/api/chat",
			bytes.NewReader(data),
		)
		if err != nil {
			errors <- err
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.Client.Do(req)
		if err != nil {
			errors <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errors <- fmt.Errorf(
				"ollama returned status: %s",
				resp.Status,
			)
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			var response ChatResponse

			if err := json.Unmarshal(
				scanner.Bytes(),
				&response,
			); err != nil {
				errors <- err
				return
			}

			if response.Message.Content != "" {
				tokens <- response.Message.Content
			}

			if response.Done {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errors <- err
			return
		}
	}()

	return tokens, errors
}
