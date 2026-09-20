package chat

type Chat struct {
	ID    int64
	Title string
}

type Message struct {
	ID      int64
	ChatID  int64
	Role    string
	Content string
}

type ChatRequest struct {
	ChatID  int64
	Role    string
	Content string
}
