package expectation

import (
	"fmt"
	"regexp"
)

type Match struct {
	pattern string
	regex *regexp.Regexp
	hasRegex bool
	name string
}

func NewMatch(pattern string, isRegex bool, name string) (Match, error) {
	regex := &regexp.Regexp{}
	var err error
	if isRegex {
		regex, err = regexp.Compile(pattern)
		if err != nil {
			return Match{}, err
		}
	}

	return Match{
		pattern: pattern,
		regex: regex,
		name: name,
		hasRegex: isRegex,
	}, nil
}

func (m Match) GetName() string {
	return m.name
}
func (m Match) GetFailure() string {
	return fmt.Sprintf("No match found for %s", m.pattern)
}

func (m Match) Assert(payload string) (isValid, isFinished bool, err error) {
	isMatched := false
	if m.hasRegex {
		isMatched = m.regex.MatchString(payload)
	} else {
		isMatched = payload == m.pattern
	}

	isValid = isMatched
	isFinished = isMatched
	return
}

