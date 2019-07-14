package tests

import (
	"log"
	"testing"
)

type TestCase func(t *testing.T, r *Runtime)

type Runtime struct {
	beforeAll  func(r *Runtime)
	beforeEach func(r *Runtime)
	afterAll   func(r *Runtime)
	afterEach  func(r *Runtime)
	testCases  []TestCase
	state      map[string]interface{}
}

func NewRuntime() *Runtime {
	return &Runtime{
		beforeAll:  func(r *Runtime) {},
		beforeEach: func(r *Runtime) {},
		afterAll:   func(r *Runtime) {},
		afterEach:  func(r *Runtime) {},
		testCases:  []TestCase{},
		state:      map[string]interface{}{},
	}
}

func (r *Runtime) SetState(key string, value interface{}) {
	r.state[key] = value
}

func (r *Runtime) GetState(key string) (interface{}, bool) {
	s, ok := r.state[key]
	return s, ok
}

func (r *Runtime) GetStateOrFail(key string) interface{} {
	s, ok := r.GetState(key)
	if !ok {
		log.Panicf("Unknown test runtime state: '%s'", key)
	}
	return s
}

func (r *Runtime) TestCase(tc TestCase) *Runtime {
	r.testCases = append(r.testCases, tc)
	return r
}

func (r *Runtime) BeforeAll(f func(*Runtime)) *Runtime {
	r.beforeAll = f
	return r
}

func (r *Runtime) BeforeEach(f func(*Runtime)) *Runtime {
	r.beforeEach = f
	return r
}

func (r *Runtime) AfterAll(f func(*Runtime)) *Runtime {
	r.afterAll = f
	return r
}

func (r *Runtime) AfterEach(f func(*Runtime)) *Runtime {
	r.afterEach = f
	return r
}

func (r *Runtime) Run(t *testing.T) {
	r.beforeAll(r)
	for _, testFunc := range r.testCases {
		r.beforeEach(r)
		testFunc(t, r)
		r.afterEach(r)
	}
	r.afterAll(r)
}
