package kafka

import "github.com/breathbath/uAssert/stream"

type StreamValidator struct {
	Address   Address
	Validator stream.StreamValidator
}

func NewStreamValidator(address Address, validator stream.StreamValidator) *StreamValidator {
	return &StreamValidator{
		Address:   address,
		Validator: validator,
	}
}

func (ka *StreamValidator) Validate(payload string) error {
	return ka.Validator.Validate(payload)
}

