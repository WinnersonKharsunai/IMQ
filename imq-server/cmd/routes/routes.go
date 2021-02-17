package routes

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/WinnersonKharsunai/IMQ/imq-server/pkg/protocol"
	"github.com/WinnersonKharsunai/IMQ/imq-server/pkg/services"
)

// Router is the reciever type
type Router struct {
	svc services.ImqServiceIF
	ptl *protocol.Protocol
}

// NewRouter is the factory funtion for the Router type
func NewRouter(svc services.ImqServiceIF, ptl *protocol.Protocol) *Router {
	return &Router{
		svc: svc,
		ptl: ptl,
	}
}

// RequestRouter handles all request and response
func (rt *Router) RequestRouter(ctx context.Context, r interface{}) interface{} {

	response := services.Response{}

	request, err := getClientRequest(r)
	if err != nil {
		response.Error = err.Error()
	}

	if err := rt.ptl.ValidateHeader(ctx, request.Header); err != nil {
		response.Error = err.Error()
	}

	switch request.Method {
	case "SendMessage":
		request, err := getSendMessageRequest(request.Body)
		if err != nil {
			response.Error = err.Error()
		}

		resp, err := rt.svc.SendMessage(ctx, request)
		if err != nil {
			response.Error = err.Error()
		}

		response.Body = *resp

	default:
		err := errors.New("unknown request method")
		response.Error = err.Error()
	}

	return response
}

func getClientRequest(body interface{}) (*services.Request, error) {

	var request *services.Request

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &request); err != nil {
		return nil, err
	}

	return request, nil
}

func getSendMessageRequest(body interface{}) (*services.SendMessageRequest, error) {

	var sendMessageRequest services.SendMessageRequest

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &sendMessageRequest); err != nil {
		return nil, err
	}

	return &sendMessageRequest, nil
}
