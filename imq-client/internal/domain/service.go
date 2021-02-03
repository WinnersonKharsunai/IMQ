package domain

import "time"

// PrepareClientRequest ...
func PrepareClientRequest(clientName, data string) Request {

	now := time.Now().UTC()

	request := Request{
		ClientName: clientName,
		Message: Message{
			CreatedAt: now.Format("2006-01-02 15:04:05"),
			ExpireAt:  now.Add(time.Duration(time.Second * 60)).Format("2006-01-02 15:04:05"),
			Data:      data,
		},
	}

	return request
}
