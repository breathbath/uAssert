package expectation

import (
	"fmt"
)

type NotMatch struct {
	match Match
}

func NewNotMatch(match Match) NotMatch {
	return NotMatch{
		match: match,
	}
}

func (m NotMatch) GetName() string {
	return m.match.GetName()
}
func (m NotMatch) GetFailure() string {
	return fmt.Sprintf("Match found for %s", m.match.pattern)
}

func (m NotMatch) Assert(payload string) (isValid, isFinished bool, err error) {
	isMatched, isFinished, err := m.match.Assert(payload)
	if isMatched {
		return false, true, nil
	}

	return true, false, nil
}

