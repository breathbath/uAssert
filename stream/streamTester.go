package stream

import (
	"context"
	"fmt"
	errs2 "github.com/breathbath/go_utils/utils/errs"
	"github.com/breathbath/go_utils/utils/io"
	"github.com/breathbath/uAssert/options"
	"github.com/breathbath/uAssert/stream/expectation"
	"sync"
	"time"
)

type StreamTester struct {
	streamFacade      StreamFacade
	expectationGroups *sync.Map
	readOptions       *sync.Map
	wg                sync.WaitGroup
}

func NewStreamTester(facade StreamFacade) *StreamTester {
	return &StreamTester{
		streamFacade:      facade,
		expectationGroups: &sync.Map{},
		readOptions:       &sync.Map{},
		wg:                sync.WaitGroup{},
	}
}

func (kct *StreamTester) AddExpectation(id string, readOptions options.Options, exp expectation.Expectation) {
	kct.readOptions.LoadOrStore(id, readOptions)

	actual, _ := kct.expectationGroups.LoadOrStore(id, &sync.Map{})

	existingExpectations := actual.(*sync.Map)
	existingExpectations.Store(exp.GetName(), exp)

	kct.expectationGroups.Store(id, existingExpectations)
	kct.wg.Add(1)
}

func (kct *StreamTester) startStreamListening(
	opts options.Options,
	outs chan string,
	errChan chan error,
) []context.CancelFunc {
	cancelFuncs := []context.CancelFunc{}

	cancelConsumptionCtx, cancelConsumptionFn := context.WithCancel(context.Background())
	kct.streamFacade.Read(
		opts,
		cancelConsumptionCtx,
		errChan,
		outs,
	)
	cancelFuncs = append(cancelFuncs, cancelConsumptionFn)

	return cancelFuncs
}

func (kct *StreamTester) startValidation(
	out chan string,
	groupName string,
	errChan chan error,
) context.CancelFunc {
	cancelValdationCtx, cancelValidationFn := context.WithCancel(context.Background())
	go func() {
	mainLoop:
		for {
			select {
			case msg := <-out:
				actualVal, ok := kct.expectationGroups.Load(groupName)
				if !ok {
					errChan <- fmt.Errorf("Unknown expectation options '%v'", groupName)
					continue mainLoop
				}

				expectationGroup, ok := actualVal.(*sync.Map)
				if !ok {
					errChan <- fmt.Errorf("Unknown type of map value '%v'", actualVal)
					continue mainLoop
				}

				keysToDelete := []string{}
				expectationGroup.Range(func(key, value interface{}) bool {
					exp, ok := value.(expectation.Expectation)
					if !ok {
						errChan <- fmt.Errorf("Unknown type of map value '%v'", value)
						return true
					}
					isSuccess, isFinished, err := exp.Assert(msg)
					if err != nil {
						errChan <- err
						return true
					}

					if isFinished {
						kct.wg.Done()
						if !isSuccess {
							errChan <- fmt.Errorf("Expectation failure: %s", exp.GetFailure())
						}
						keysToDelete = append(keysToDelete, exp.GetName())
					}

					return true
				})

				if len(keysToDelete) > 0 {
					for _, keyToDelete := range keysToDelete {
						expectationGroup.Delete(keyToDelete)
					}
				}
				kct.expectationGroups.Store(groupName, expectationGroup)
			case <-cancelValdationCtx.Done():
				return
			}
		}

	}()
	return cancelValidationFn
}

func (kct *StreamTester) startCollectingErrors(errChan chan error, errCont *errs2.ErrorContainer) chan bool {
	doneChan := make(chan bool)
	go func() {
		defer close(doneChan)
		for err := range errChan {
			errCont.AddError(err)
		}
	}()

	return doneChan
}

func (kct *StreamTester) StartTesting(timeout time.Duration) *errs2.ErrorContainer {
	cancelFuncs := []context.CancelFunc{}
	errChan := make(chan error)
	errCont := errs2.NewErrorContainer()

	doneCollectingErrsChan := kct.startCollectingErrors(errChan, errCont)

	outs := make(chan string)
	kct.expectationGroups.Range(func(key, value interface{}) bool {
		optKey := key.(string)
		opts, _ := kct.readOptions.Load(optKey)

		readOptions := opts.(options.Options)
		curCancelFuncs := kct.startStreamListening(readOptions, outs, errChan)

		cancelFuncs = append(cancelFuncs, curCancelFuncs...)

		curCancelFunc := kct.startValidation(outs, optKey, errChan)
		cancelFuncs = append(cancelFuncs, curCancelFunc)

		return true
	})

	wgChan := make(chan struct{})
	go func() {
		defer close(wgChan)
		kct.wg.Wait()
	}()

	defer func() {
		close(errChan)
		<-doneCollectingErrsChan
	}()

	select {
	case <-wgChan:
		io.OutputInfo("", "All requirements are satisfied")
	case <-time.After(timeout):
		errCont.AddErrorF("Timeout after %v", timeout)
	}

	for _, cancelFn := range cancelFuncs {
		cancelFn()
	}

	return errCont
}
