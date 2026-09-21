package chat

type Chat struct {
	ID        int64
	Title     string
	CreatedAt string
	UpdatedAt string
}

type Message struct {
	ID        int64
	ChatID    int64
	Role      string
	Content   string
	CreatedAt string
}
