package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	Client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 3 * time.Minute,
		},
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatResponse struct {
	Model   string      `json:"model"`
	Message ChatMessage `json:"message"`
	Done    bool        `json:"done"`
}

type Model struct {
	Name string `json:"name"`
}

type ModelResponse struct {
	Models []Model `json:"models"`
}

func (c *Client) GetModels() ([]Model, error) {
	resp, err := c.Client.Get(c.BaseURL + "/api/tags")

	if err != nil {
		return nil, fmt.Errorf("Get Model Error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("Ollama Returned %s:%s", resp.Status, string(body))
	}

	var result ModelResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Decode Models Error: %w", err)
	}

	return result.Models, nil
}

func (c *Client) Chat(model string, messages []ChatMessage) (<-chan string, <-chan error) {
	tokenChan := make(chan string)
	errorChan := make(chan error, 1)

	go func() {
		defer close(tokenChan)
		defer close(errorChan)

		requestBody := ChatRequest{
			Model:    model,
			Messages: messages,
			Stream:   true,
		}

		data, err := json.Marshal(requestBody)

		if err != nil {
			errorChan <- fmt.Errorf(
				"Marshal Request Error: %w",
				err,
			)
			return
		}

		resp, err := c.Client.Post(
			c.BaseURL+"/api/chat",
			"application/json",
			bytes.NewReader(data),
		)
		if err != nil {
			errorChan <- fmt.Errorf("Send Chat Request Error: %w", err)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)

			errorChan <- fmt.Errorf("Ollama Returned %s:%s", resp.Status, string(body))
			return
		}

		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			var response ChatResponse

			if err := json.Unmarshal(
				scanner.Bytes(),
				&response,
			); err != nil {
				errorChan <- fmt.Errorf("Decode Stream Error: %w", err)
				return
			}

			if response.Message.Content != "" {
				tokenChan <- response.Message.Content
			}

			if response.Done {
				return
			}

			if err := scanner.Err(); err != nil {
				errorChan <- fmt.Errorf("Read Stream Error: %w", err)
			}
		}
	}()

	return tokenChan, errorChan

}
