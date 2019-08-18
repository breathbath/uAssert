package expectation

import (
	"fmt"
	"strings"
	"sync"
)

type Sequence struct {
	expectationsSequence []Expectation
	name string
	curIndex int
	mux sync.Mutex
}

func NewSequence(expectationsSequence []Expectation, name string) *Sequence {
	return &Sequence{expectationsSequence: expectationsSequence, name: name, curIndex: 0, mux: sync.Mutex{}}
}

func (s *Sequence) GetName() string {
	return s.name
}
func (s *Sequence) GetFailure() string {
	sequence := make([]string, len(s.expectationsSequence))
	for k, s := range s.expectationsSequence {
		sequence[k] = s.GetName()
	}

	return fmt.Sprintf("Sequence of expectations '%s' was not discovered ", strings.Join(sequence, ", "))
}

func (s *Sequence) Assert(payload string) (isValid, isFinished bool, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.curIndex >= len(s.expectationsSequence) {
		return true, true, nil
	}

	curSequence := s.expectationsSequence[s.curIndex]
	curIsMatched, curIsFinished, curErr := curSequence.Assert(payload)
	if curIsMatched && curIsFinished {
		s.curIndex++
	}
	err = curErr

	if err == nil && s.curIndex >= len(s.expectationsSequence) {
		return true, true, nil
	}

	return
}