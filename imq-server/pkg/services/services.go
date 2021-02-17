package services

import "context"

// ImqServiceIF defines the methods available for IMQ Service
type ImqServiceIF interface {
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
	Data      string `json:"data"`
	CreatedAt string `json:"createdAt"`
	ExpireAt  string `json:"expireAt"`
}

// Request holds the request sent by client
type Request struct {
	Header interface{} `json:"header"`
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
}

// Response holds the response to be sent to client
type Response struct {
	Error string      `json:"error" validate:"omitempty"`
	Body  interface{} `json:"body" validate:"omitempty"`
}
