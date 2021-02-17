package domain

// Message type holds client's message
type Message struct {
	ClientName string `json:"clientName"`
	Data       string `json:"data"`
	CreatedAt  string `json:"createdAt"`
	ExpireAt   string `json:"expireAt"`
}

// Request type holds client's request
type Request struct {
	Header Header      `json:"header"`
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
}

// Header ...
type Header struct {
	Version string `json:"version"`
}

// SendMessageRequest ...
type SendMessageRequest struct {
	ClientName string `json:"clientName"`
	Data       string `json:"data"`
	CreatedAt  string `json:"createdAt"`
	ExpireAt   string `json:"expireAt"`
}
