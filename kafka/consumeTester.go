package kafka

import (
	"context"
	"fmt"
	errs2 "github.com/breathbath/go_utils/utils/errs"
	"github.com/breathbath/go_utils/utils/io"
	"github.com/breathbath/uAssert/stream"
	"sync"
	"time"
)

type ConsumeTester struct {
	kafkaFacade     Facade
	validatorGroups map[Address][]*StreamValidator
}

const Consumers_Per_Address_Count = 5

func NewKafkaConsumerTester(kafkaFacade Facade) *ConsumeTester {
	return &ConsumeTester{
		kafkaFacade:     kafkaFacade,
		validatorGroups: map[Address][]*StreamValidator{},
	}
}

func (kct *ConsumeTester) AddValidator(address Address, validator stream.Validator) {
	_, ok := kct.validatorGroups[address]
	if !ok {
		kct.validatorGroups[address] = []*StreamValidator{}
	}
	kct.validatorGroups[address] = append(kct.validatorGroups[address], NewStreamValidator(address, validator))
}

func (kct *ConsumeTester) startConsumption(
	kafkaAddr Address,
	cancelFuncs *[]context.CancelFunc,
	errChans *[]chan error,
	msgsStreams *[]<-chan string,
) {
	for i := 0; i < Consumers_Per_Address_Count; i++ {
		cancelConsumptionCtx, cancelConsumptionFn := context.WithCancel(context.Background())
		msgChan, errChan := kct.kafkaFacade.Read(
			ConsumerOptions{kafkaAddr, 5},
			cancelConsumptionCtx,
		)
		*cancelFuncs = append(*cancelFuncs, cancelConsumptionFn)
		*errChans = append(*errChans, errChan)
		*msgsStreams = append(*msgsStreams, msgChan)
	}
}

func (kct *ConsumeTester) startValidation(
	msgsStreams []<-chan string,
	validationGroup []*StreamValidator,
	cancelFuncs *[]context.CancelFunc,
	errChans *[]chan error,
) {
	cancelValdationCtx, cancelValidationFn := context.WithCancel(context.Background())
	msgValidateErrChan := make(chan error)
	go func() {
		fannedInMsgStream := stream.MergeMessageStreams(msgsStreams...)
		for {
			select {
			case msg := <-fannedInMsgStream:
				for _, validator := range validationGroup {
					err := validator.Validate(msg)
					if err != nil {
						msgValidateErrChan <- err
					}
				}
			case <-cancelValdationCtx.Done():
				return
			}
		}

	}()
	*cancelFuncs = append(*cancelFuncs, cancelValidationFn)
	*errChans = append(*errChans, msgValidateErrChan)
}

func (kct *ConsumeTester) startFinishedValidatorsTracker(
	validationGroup []*StreamValidator,
	wg *sync.WaitGroup,
	cancelFuncs *[]context.CancelFunc,
) {
	cancelValdationCheckCtx, cancelValidationCheckFn := context.WithCancel(context.Background())
	go func() {
		for {
			allAreFinished := true
			for _, validator := range validationGroup {
				isFinished := validator.Validator.IsFinished()
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
	*cancelFuncs = append(*cancelFuncs, cancelValidationCheckFn)
}

func (kct *ConsumeTester) collectErrors(
	errChans []chan error,
	errs *errs2.ErrorContainer,
	cancelFuncs *[]context.CancelFunc,
) {
	cancelCollectErrsCtx, cancelCollectErrsFn := context.WithCancel(context.Background())
	stream.CollectErrors(errs, cancelCollectErrsCtx, errChans...)
	*cancelFuncs = append(*cancelFuncs, cancelCollectErrsFn)
}

func (kct *ConsumeTester) StartTesting(timeout time.Duration) error {
	errs := errs2.NewErrorContainer()

	wg := sync.WaitGroup{}
	wg.Add(len(kct.validatorGroups))

	cancelFuncs := []context.CancelFunc{}
	errChans := make([]chan error, 0)

	for kafkaAddr, validationGroup := range kct.validatorGroups {
		msgsStreams := make([]<-chan string, Consumers_Per_Address_Count)

		kct.startConsumption(kafkaAddr, &cancelFuncs, &errChans, &msgsStreams)

		kct.startValidation(msgsStreams, validationGroup, &cancelFuncs, &errChans)

		kct.startFinishedValidatorsTracker(validationGroup, &wg, &cancelFuncs)

		kct.collectErrors(errChans, errs, &cancelFuncs)
	}

	wgChan := make(chan struct{})
	go func() {
		defer close(wgChan)
		wg.Wait()
	}()

	select {
	case <-wgChan:
		io.OutputInfo("", "All requirements are satisfied")
	case <-time.After(timeout):
		errs.AddError(fmt.Errorf("Timeout after %v", timeout))
		io.OutputInfo("", "Testing process timeout")
	}

	for _, cancelFn := range cancelFuncs {
		cancelFn()
	}

	return errs.Result(" ")
}
