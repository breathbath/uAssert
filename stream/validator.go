package stream

type Validator interface {
	Validate(streamItem string) (err error)
	GetValidationErrors() []error
	IsFinished() bool
}