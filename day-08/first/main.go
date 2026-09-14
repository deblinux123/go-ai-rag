package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func main() {
	requestBody := ChatRequest{
		Model: "qwen2.5:3b",
		Messages: []Message{
			{
				Role:    "user",
				Content: "What is go programming language? tell me in persian",
			},
		},
		Stream: false,
	}

	jsonBody, err := json.Marshal(requestBody)

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		ollamaURL+"/api/chat",
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		panic(fmt.Sprintf(
			"Ollama returned %s: %s",
			resp.Status,
			string(body),
		))
	}

	var result ChatResponse

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		panic(err)
	}

	fmt.Println("Model:", result.Model)
	fmt.Println("Role:", result.Message.Role)
	fmt.Println("Response:")
	fmt.Println(result.Message.Content)
	fmt.Println("Done:", result.Done)
}
