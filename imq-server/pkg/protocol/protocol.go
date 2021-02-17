package protocol

import (
	"context"
	"encoding/json"
	"errors"
)

// Header holds the header details
type Header struct {
	Version string `json:"version"`
}

// Protocol is the reciever type
type Protocol struct{}

// NewProtocol is the factory function for the Protocol
func NewProtocol() *Protocol {
	return &Protocol{}
}

// ValidateHeader validates incoming header
func (p *Protocol) ValidateHeader(ctx context.Context, header interface{}) error {

	var hdr Header

	b, err := json.Marshal(header)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &hdr); err != nil {
		return err
	}

	if hdr.Version != "1.0" {
		return errors.New("given version not implemented")
	}

	return nil
}
