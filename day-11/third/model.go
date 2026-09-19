package main

type Chat struct {
	ID    int
	Title string
}

type ChatMessage struct {
	ID      int
	ChatID  int
	Role    string
	Content string
}
