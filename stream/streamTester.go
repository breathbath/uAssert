package stream

import (
	"context"
	errs2 "github.com/breathbath/go_utils/utils/errs"
	"github.com/breathbath/go_utils/utils/io"
	"github.com/breathbath/uAssert/options"
	"github.com/breathbath/uAssert/stream/validate"
	"strings"
	"sync"
	"time"
)

type StreamTester struct {
	streamFacade    StreamFacade
	validatorGroups map[string][]validate.Validator
}

const Consumers_Per_Address_Count = 1

func NewStreamTester(facade StreamFacade) *StreamTester {
	return &StreamTester{
		streamFacade:    facade,
		validatorGroups: map[string][]validate.Validator{},
	}
}

func (kct *StreamTester) AddValidator(id string, validator validate.Validator) {
	_, ok := kct.validatorGroups[id]
	if !ok {
		kct.validatorGroups[id] = []validate.Validator{}
	}
	kct.validatorGroups[id] = append(kct.validatorGroups[id], validator)
}

func (kct *StreamTester) startStreamListening(
	opts options.Options,
	outs chan string,
	errChan chan error,
) []context.CancelFunc {
	cancelFuncs := []context.CancelFunc{}
	for i := 0; i < Consumers_Per_Address_Count; i++ {
		cancelConsumptionCtx, cancelConsumptionFn := context.WithCancel(context.Background())
		kct.streamFacade.Read(
			opts,
			cancelConsumptionCtx,
			errChan,
			outs,
		)
		cancelFuncs = append(cancelFuncs, cancelConsumptionFn)
	}

	return cancelFuncs
}

func (kct *StreamTester) startValidation(
	out chan string,
	validationGroup []validate.Validator,
	errChan chan error,
) context.CancelFunc {
	cancelValdationCtx, cancelValidationFn := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case msg := <-out:
				for _, validator := range validationGroup {
					err := validator.Validate(msg)
					if err != nil {
						errChan <- err
					}
				}
			case <-cancelValdationCtx.Done():
				return
			}
		}

	}()
	return cancelValidationFn
}

func (kct *StreamTester) startFinishedValidatorsTracker(
	validationGroup []validate.Validator,
	wg *sync.WaitGroup,
) context.CancelFunc {
	cancelValdationCheckCtx, cancelValidationCheckFn := context.WithCancel(context.Background())
	go func() {
		for {
			allAreFinished := true
			for _, validator := range validationGroup {
				isFinished := validator.IsFinished()
				if !isFinished {
					allAreFinished = false
				}
			}
			if allAreFinished {
				wg.Done()
				return
			}

			select {
			case <-time.After(time.Second * 1):
				continue
			case <-cancelValdationCheckCtx.Done():
				return
			}
		}
	}()
	return cancelValidationCheckFn
}

func (kct *StreamTester) startCollectingErrors(errChan chan error, errCont *errs2.ErrorContainer) {
	go func() {
		for err := range errChan {
			errCont.AddError(err)
		}
	}()
}

func (kct *StreamTester) StartTesting(opts options.Options, timeout time.Duration) *errs2.ErrorContainer {
	cancelFuncs := []context.CancelFunc{}
	errChan := make(chan error)
	errCont := errs2.NewErrorContainer()

	kct.startCollectingErrors(errChan, errCont)

	wg := sync.WaitGroup{}
	wg.Add(len(kct.validatorGroups))

	outs := make(chan string)
	for _, validationGroup := range kct.validatorGroups {
		curCancelFuncs := kct.startStreamListening(opts, outs, errChan)
		cancelFuncs = append(cancelFuncs, curCancelFuncs...)

		curCancelFunc := kct.startValidation(outs, validationGroup, errChan)
		cancelFuncs = append(cancelFuncs, curCancelFunc)

		curCancelFunc = kct.startFinishedValidatorsTracker(validationGroup, &wg)
		cancelFuncs = append(cancelFuncs, curCancelFunc)
	}

	wgChan := make(chan struct{})
	go func() {
		defer close(wgChan)
		wg.Wait()
	}()

	defer func() {
		kct.collectValidationErrors(errCont)
		close(errChan)
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

func (kct *StreamTester) collectValidationErrors(errCont *errs2.ErrorContainer) {
	for _, validationGroup := range kct.validatorGroups {
		for _, validator := range validationGroup {
			currErrs := validator.GetValidationErrors()
			if len(currErrs) == 0 {
				continue
			}

			validatorErrors := make([]string, 0, len(currErrs))
			for _, curErr := range currErrs {
				validatorErrors = append(validatorErrors, curErr.Error())
			}

			errCont.AddErrorF("Validator '%s' has failed: %s", validator.GetName(), strings.Join(validatorErrors, "; "))
		}
	}
}
