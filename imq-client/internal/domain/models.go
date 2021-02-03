package domain

// Message type holds client's message
type Message struct {
	Data      string `json:"data"`
	CreatedAt string `json:"createdAt"`
	ExpireAt  string `json:"expireAt"`
}

// Request type holds client's request
type Request struct {
	ClientName string  `json:"clientName"`
	Message    Message `json:"message"`
}
