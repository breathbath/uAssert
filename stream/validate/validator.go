package validate

type Validator interface {
	Validate(streamItem string) (err error)
	GetValidationErrors() []error
	IsFinished() bool
	GetName() string
}