package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	Client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Client: &http.Client{
			Timeout: 0,
		},
	}
}

func (c *Client) SetBaseURL(baseURL string) {
	c.BaseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	Stream    bool          `json:"stream"`
	Options   ChatOptions   `json:"options,omitempty"`
	KeepAlive string        `json:"keep_alive,omitempty"`
}

type ChatOptions struct {
	Temperature float64 `json:"temperature"`
}

type ChatResponse struct {
	Model      string      `json:"model"`
	Message    ChatMessage `json:"message"`
	Done       bool        `json:"done"`
	DoneReason string      `json:"done_reason,omitempty"`
}

type Model struct {
	Name string `json:"name"`
}

type ModelsResponse struct {
	Models []Model `json:"models"`
}

func (c *Client) GetModels() ([]Model, error) {
	if c.BaseURL == "" {
		return nil, fmt.Errorf("ollama URL is empty")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(c.BaseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("get models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var result ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode models: %w", err)
	}

	return result.Models, nil
}

func (c *Client) Chat(
	model string,
	messages []ChatMessage,
	temperature float64,
) (<-chan string, <-chan error) {
	tokenChan := make(chan string)
	errorChan := make(chan error, 1)

	go func() {
		defer close(tokenChan)
		defer close(errorChan)

		requestBody := ChatRequest{
			Model:    model,
			Messages: messages,
			Stream:   true,
			Options: ChatOptions{
				Temperature: temperature,
			},
			KeepAlive: "5m",
		}

		data, err := json.Marshal(requestBody)
		if err != nil {
			errorChan <- fmt.Errorf("marshal request: %w", err)
			return
		}

		req, err := http.NewRequest(
			http.MethodPost,
			c.BaseURL+"/api/chat",
			bytes.NewReader(data),
		)
		if err != nil {
			errorChan <- fmt.Errorf("create chat request: %w", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.Client.Do(req)
		if err != nil {
			errorChan <- fmt.Errorf("send chat request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errorChan <- fmt.Errorf(
				"ollama returned %s: %s",
				resp.Status,
				strings.TrimSpace(string(body)),
			)
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)

		for scanner.Scan() {
			var response ChatResponse

			if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
				errorChan <- fmt.Errorf("decode stream: %w", err)
				return
			}

			if response.Message.Content != "" {
				tokenChan <- response.Message.Content
			}

			if response.Done {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errorChan <- fmt.Errorf("read stream: %w", err)
		}
	}()

	return tokenChan, errorChan
}
