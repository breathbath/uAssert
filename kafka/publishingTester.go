package kafka

import (
	"context"
	"github.com/breathbath/go_utils/utils/io"
	"time"
)

type PublishingPayload struct {
	PublishingDelay time.Duration
	Body            string
}

type PublishingTester struct {
	kafkaFacade Facade
}

func NewPublishingTester(kafkaFacade Facade) PublishingTester {
	return PublishingTester{kafkaFacade: kafkaFacade}
}

func (pt PublishingTester) PublishPayloadsAsync(
	opts ProducerOptions,
	payloads []PublishingPayload,
	errs chan error,
) {
	defer close(errs)
	for _, payload := range payloads {
		go func(payload PublishingPayload) {
			io.OutputInfo("", "Got to publish %s", payload.Body)
			time.Sleep(payload.PublishingDelay)

			err := pt.kafkaFacade.Publish(payload.Body, opts, context.Background())
			if err != nil {
				errs <- err
			} else {
				io.OutputInfo("", "Published %s", payload.Body)
			}
		}(payload)
	}
}
