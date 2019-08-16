package kafka

import (
	"context"
	"fmt"
	"github.com/breathbath/go_utils/utils/env"
	"github.com/breathbath/go_utils/utils/errs"
	"github.com/breathbath/go_utils/utils/io"
	"github.com/breathbath/uAssert/stream"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TOPIC_TO_TEST = "test_kafka"

var kafkaFacade Facade
var connector *Connector

func setup() {
	connStr := env.ReadEnvOrFail("KAFKA_CONN_STR")
	connector = NewConnector(connStr)

	var err error
	kafkaFacade = NewFacade(connector, connStr)
	errs.FailOnError(err)

	err = kafkaFacade.CreateTopics(context.Background(), TOPIC_TO_TEST)
	errs.FailOnError(err)
}

func cleanup() {
	fmt.Println("cleaning topics")
	err := kafkaFacade.CleanTopics(TOPIC_TO_TEST)
	if err != nil {
		io.OutputError(err, "", "Failed to cleanup topic %s", TOPIC_TO_TEST)
	}

	err = connector.Destroy()
	if err != nil {
		io.OutputError(err, "", "Failed to destroy kafka connector")
	}
}

func TestFacade(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	setup()
	defer cleanup()
	t.Run("testProducing", testProducing)
	t.Run("testConsumption", testConsumption)
}

func startErrorsListener(ctx context.Context, errsCont *errs.ErrorContainer) chan error {
	errsChan := make(chan error)
	stream.CollectErrors(errsCont, ctx, errsChan)

	return errsChan
}

func testConsumption(t *testing.T) {
	errsCont := errs.NewErrorContainer()
	cancelCollectErrsCtx, cancelCollectErrsFn := context.WithCancel(context.Background())
	errsChan := startErrorsListener(cancelCollectErrsCtx, errsCont)
	defer cancelCollectErrsFn()

	opts := ProducerOptions{
		addr:             Address{Topic: TOPIC_TO_TEST},
		writeDeadLineSec: 1,
	}

	publisherTester := NewPublishingTester(kafkaFacade)
	publisherTester.PublishPayloadsAsync(
		opts,
		[]PublishingPayload{
			{
				PublishingDelay: time.Duration(0),
				Body:            "consumption lala",
			},
			{
				PublishingDelay: time.Duration(20),
				Body:            "2 consumption",
			},
		},
		errsChan,
	)

	consumerTester := NewKafkaConsumeTester(kafkaFacade)

	exactMatch := stream.NewExactMatchAssertion("consumption lala", )
	consumerTester.AddValidator(Address{Topic: TOPIC_TO_TEST}, exactMatch)

	regexValidator := stream.NewRegexValidator(`^\d+`)
	consumerTester.AddValidator(Address{Topic: TOPIC_TO_TEST}, regexValidator)

	consumerTester.StartTesting(time.Second*10, errsChan)
	assert.NoError(t, errsCont.Result(" "))
}

func testProducing(t *testing.T) {
	opts := ProducerOptions{
		addr:             Address{Topic: TOPIC_TO_TEST},
		writeDeadLineSec: 1,
	}
	err := kafkaFacade.Publish("some payload", opts, context.Background())
	assert.NoError(t, err)
}
