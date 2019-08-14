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

func testConsumption(t *testing.T) {
	consumerTester := NewKafkaConsumerTester(kafkaFacade)

	exachtMatch := stream.NewExactMatchAssertion("some payload")
	kafkaExactMatch := NewStreamValidator(Address{Topic: TOPIC_TO_TEST}, exachtMatch)
	consumerTester.AddValidator(kafkaExactMatch)

	err := consumerTester.StartTesting(time.Second * 10)
	assert.NoError(t, err)
}

func testProducing(t *testing.T) {
	opts := ProducerOptions{
		addr:             Address{Topic: TOPIC_TO_TEST},
		writeDeadLineSec: 1,
	}
	err := kafkaFacade.Publish("some payload", opts, context.Background())
	assert.NoError(t, err)
}