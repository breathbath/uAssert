package stream

import (
"fmt"
"regexp"
)

type RegexValidator struct {
	regex          *regexp.Regexp
	validationErrors []error
}

func NewRegexValidator(regex string) *RegexValidator {
	return &RegexValidator{
		regex: regexp.MustCompile(regex),
		validationErrors: []error{
			fmt.Errorf("Nothing matched with the regex '%s'", regex),
		},
	}
}

func (rv *RegexValidator) Validate(payload string) (err error) {
	if rv.regex.MatchString(payload) {
		rv.validationErrors = []error{}
	}

	return nil
}

func (rv *RegexValidator) GetValidationErrors() []error {
	return rv.validationErrors
}

func (rv *RegexValidator) IsFinished() bool {
	return len(rv.validationErrors) == 0
}


