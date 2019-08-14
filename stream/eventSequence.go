package stream

import (
	"fmt"
)

type EventSequenceValidator struct {
	eventsValidators [] Validator
	validationErrors []error
	nextValidator int
}

func NewEventSequenceValidator(eventsValidators [] Validator) (*EventSequenceValidator, error) {
	var err error
	if len(eventsValidators) == 0 {
		err = fmt.Errorf("No sequence expectations provided")
	}
	return &EventSequenceValidator{
		eventsValidators: [] Validator{},
		validationErrors: []error{
			fmt.Errorf("Nothing matched to the defined sequence"),
		},
		nextValidator: 0,
	}, err
}

func (esv *EventSequenceValidator) Validate(payload string) (err error) {
	if esv.nextValidator >= len(esv.eventsValidators) {
		return nil
	}
	nextValidator := esv.eventsValidators[esv.nextValidator]
	err = nextValidator.Validate(payload)
	if err != nil {
		return err
	}

	if nextValidator.IsFinished() {
		esv.nextValidator = esv.nextValidator + 1
	}

	return nil
}

func (esv *EventSequenceValidator) GetValidationErrors() []error {
	errs := []error{}
	for _, esv := range esv.eventsValidators {
		errs = append(errs, esv.GetValidationErrors()...)
	}
	return errs
}

func (esv *EventSequenceValidator) IsFinished() bool {
	return esv.nextValidator >= len(esv.eventsValidators)
}


