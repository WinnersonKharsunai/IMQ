package domain

import "time"

// PrepareClientRequest ...
func PrepareClientRequest(clientName, data string) Request {

	now := time.Now().UTC()

	sm := SendMessageRequest{
		ClientName: clientName,
		CreatedAt:  now.Format("2006-01-02 15:04:05"),
		ExpireAt:   now.Add(time.Duration(time.Second * 60)).Format("2006-01-02 15:04:05"),
		Data:       data,
	}

	request := Request{
		Header: Header{
			Version: "1.0",
		},
		Method: "SendMessage",
		Body:   sm,
	}

	return request
}
