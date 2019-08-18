package validate

import (
	"fmt"
	"strings"
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
		eventsValidators: eventsValidators,
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

func (esv *EventSequenceValidator) GetName() string {
	names := []string{}
	for _, v := range esv.eventsValidators {
		names = append(names, v.GetName())
	}

	return fmt.Sprintf("Expecting event sequence: %s", strings.Join(names, ", "))
}

func (esv *EventSequenceValidator) GetValidationErrors() []error {
	return esv.validationErrors
}

func (esv *EventSequenceValidator) IsFinished() bool {
	sequenceIsFound := esv.nextValidator >= len(esv.eventsValidators)
	if sequenceIsFound {
		esv.validationErrors = []error{}
	}

	return sequenceIsFound
}


