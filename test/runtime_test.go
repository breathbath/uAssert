package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutionOrder(t *testing.T) {
	actualEvents := []string{}

	r := NewRuntime()

	r.TestCase(func(t *testing.T, r *Runtime) {
		actualEvents = append(actualEvents, "test_case_1")
	})
	r.BeforeEach(func(r *Runtime) {
		actualEvents = append(actualEvents, "before_each")
	})
	r.AfterEach(func(r *Runtime) {
		actualEvents = append(actualEvents, "after_each")
	})
	r.AfterAll(func(r *Runtime) {
		actualEvents = append(actualEvents, "after_all")
	})
	r.TestCase(func(t *testing.T, r *Runtime) {
		actualEvents = append(actualEvents, "test_case_2")
	})
	r.BeforeAll(func(r *Runtime) {
		actualEvents = append(actualEvents, "before_all")
	})

	r.Run(t)

	expectedEventsOrder := []string{
		"before_all",
		"before_each",
		"test_case_1",
		"after_each",
		"before_each",
		"test_case_2",
		"after_each",
		"after_all",
	}

	assert.Equal(t, expectedEventsOrder, actualEvents)
}

func TestStatefulExecution(t *testing.T) {
	r := NewRuntime()

	r.TestCase(func(t *testing.T, r *Runtime) {
		r.SetState("TestCase1", "TestCase1")
	})
	r.BeforeEach(func(r *Runtime) {
		r.SetState("BeforeEach", "BeforeEach")
	})
	r.AfterEach(func(r *Runtime) {
		r.SetState("AfterEach", "AfterEach")
	})
	r.AfterAll(func(r *Runtime) {
		r.SetState("AfterAll", "AfterAll")
	})
	r.TestCase(func(t *testing.T, r *Runtime) {
		r.SetState("TestCase2", "TestCase2")
	})
	r.BeforeAll(func(r *Runtime) {
		r.SetState("BeforeAll", "BeforeAll")
	})

	r.Run(t)

	expectedStates := []string{
		"TestCase1",
		"BeforeEach",
		"AfterEach",
		"AfterAll",
		"TestCase2",
		"BeforeAll",
	}

	for _, expectedState := range expectedStates {
		actualState, found := r.GetState(expectedState)

		assert.True(t, found)
		assert.Equal(t, expectedState, actualState)

		actualState = r.GetStateOrFail(expectedState)
		assert.Equal(t, expectedState, actualState)
	}
}

func TestMissingState(t *testing.T) {
	r := NewRuntime()
	r.Run(t)
	actualState, found := r.GetState("some state")
	assert.False(t, found)
	assert.Nil(t, actualState)
}

func TestRemovingState(t *testing.T) {
	r := NewRuntime()
	r.SetState("someKey", "someVal")
	r.RemoveState("someKey")
	actualState, found := r.GetState("someKey")
	assert.False(t, found)
	assert.Nil(t, actualState)
}

func TestMissingStateFailure(t *testing.T) {
	r := NewRuntime()
	r.Run(t)

	assert.PanicsWithValue(
		t,
		"Unknown test runtime state: 'some state'",
		func() {
			r.GetStateOrFail("some state")
		},
	)
}

func TestEmptySetup(t *testing.T) {
	r := NewRuntime()
	r.Run(t)

	r.TestCase(func(t *testing.T, r *Runtime) {})
	r.Run(t)
}
