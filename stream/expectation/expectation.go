package expectation

type Expectation interface {
	GetName() string
	GetFailure() string
	Assert(payload string) (isValid, isFinished bool, err error)
}
