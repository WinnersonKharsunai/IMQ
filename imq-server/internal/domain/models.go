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

// Response type holds client's response
type Response struct {
	Name      string `json:"name"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}
