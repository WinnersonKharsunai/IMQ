package protocol

import (
	"errors"
)

const (
	_version    = "1.0"
	sendMessage = "SendMessage"
)

type Validator struct{}

// Header validaites the request header
func (v *Validator) Header(version, method string) error {

	if version != _version {
		return errors.New("requested version not implemented")
	}

	switch method {
	case sendMessage:
		return nil
	default:
		return errors.New("method not implemented")
	}
}
