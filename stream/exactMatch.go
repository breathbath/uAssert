package stream

import (
	"fmt"
)

type ExactMatchValidator struct {
	pattern          string
	validationErrors []error
}

func NewExactMatchAssertion(pattern string) *ExactMatchValidator {
	return &ExactMatchValidator{
		pattern: pattern,
		validationErrors: []error{
			fmt.Errorf("No exact match found for '%s'", pattern),
		},
	}
}

func (ema *ExactMatchValidator) Validate(payload string) (err error) {
	if payload == ema.pattern {
		ema.validationErrors = []error{}
	}

	return nil
}

func (ema *ExactMatchValidator) GetValidationErrors() []error {
	return ema.validationErrors
}

func (ema *ExactMatchValidator) IsFinished() bool {
	return len(ema.validationErrors) == 0
}
