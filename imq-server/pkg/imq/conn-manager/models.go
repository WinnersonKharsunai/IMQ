package server

import "context"

// Request holds the request send by client
type Request struct {
	Header Header      `json:"header"`
	Body   interface{} `json:"body"`
}

// Header holds the request header
type Header struct {
	Version string `json:"version"`
	Method  string `json:"method"`
}

// ImqServerServiceIF ...
type ImqServerServiceIF interface {
	SendMessage(ctx context.Context, in *SendMessageRequest) (*SendMessageResponse, error)
}

// SendMessageRequest holds send message request data
type SendMessageRequest struct {
	ClientName string `json:"clientName"`
	Data       string `json:"data"`
	CreatedAt  string `json:"createdAt"`
	ExpireAt   string `json:"expireAt"`
}

// SendMessageResponse holds the send message response data
type SendMessageResponse struct {
	_ struct{}
}
