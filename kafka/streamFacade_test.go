package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/breathbath/go_utils/utils/env"
	"github.com/breathbath/go_utils/utils/errs"
	"github.com/breathbath/go_utils/utils/io"
	"github.com/breathbath/uAssert/options"
	"github.com/breathbath/uAssert/stream"
	"github.com/breathbath/uAssert/stream/expectation"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TOPIC_TO_TEST = "test_kafka"

var kafkaFacade *StreamFacade

func setup() {
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_2_0

	connStr := env.ReadEnvOrFail("KAFKA_CONN_STR")

	var err error
	kafkaFacade, err = NewStreamFacade(connStr)
	errs.FailOnError(err)

	err = kafkaFacade.PrepareStream(options.Options{"topic": TOPIC_TO_TEST})
	errs.FailOnError(err)
}

func cleanup() {
	err := kafkaFacade.Close(options.Options{"topic": TOPIC_TO_TEST})
	if err != nil {
		io.OutputError(err, "", "Failed to destroy kafka facade")
	}
}

func TestFacade(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	setup()
	defer cleanup()
	t.Run("testKafkaStream", testKafkaStream)
}

func testKafkaStream(t *testing.T) {
	opts := options.Options{
		"topic":            TOPIC_TO_TEST,
		"partition":        0,
		"writeDeadLineSec": 1,
	}

	err := kafkaFacade.PublishMany(opts, []string{"consumption lala", "2 consumption"})
	errs.FailOnError(err)

	streamTester := stream.NewStreamTester(kafkaFacade)

	readOptions := options.Options{
		"readDeadlineSec": 1,
		"topic":           TOPIC_TO_TEST,
		"partition":       0,
	}

	exactMatch, err := expectation.NewMatch(
		"consumption lala",
		false,
		"Exact match 'consumption lala'",
	)
	errs.FailOnError(err)
	streamTester.AddExpectation(TOPIC_TO_TEST, readOptions, exactMatch)

	regexValidator, err := expectation.NewMatch(
		`^\d+`,
		true,
		"Message starting with at least one number",
	)
	streamTester.AddExpectation(TOPIC_TO_TEST, readOptions, regexValidator)

	eventSequence := expectation.NewSequence(
		[]expectation.Expectation{
			exactMatch,
			regexValidator,
		},
		"First 'consumption lala' second a message with at least a number in the beginning",
	)
	streamTester.AddExpectation(TOPIC_TO_TEST, readOptions, eventSequence)
	errs.FailOnError(err)

	errsCont := streamTester.StartTesting(time.Second*10)
	io.OutputInfo("", "Draining errs")
	assert.NoError(t, errsCont.Result("\n"))
}
