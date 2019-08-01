package simulation

import (
	"context"
	"github.com/breathbath/go_utils/utils/env"
	"github.com/breathbath/uAssert/options"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKafkaIsAvailableForProducer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	topic := "test_producer"
	connector := NewKafkaConnector()
	producer := NewKafkaProducer(connector)
	pubOptions := options.Options{
		"topic":               topic,
		"conn":                env.ReadEnvOrFail("KAFKA_CONN_STR"),
		"write_dead_line_sec": 10,
		"max_attempts": 3,
	}

	err := producer.Publish("some payload", pubOptions, context.Background())
	assert.NoError(t, err)

	err = cleanTopics(pubOptions, connector)
	assert.NoError(t, err)

	err = connector.Destroy()
	assert.NoError(t, err)
}

func cleanTopics(connOpts options.Options, connector *KafkaConnector) error {
	cl := NewKafkaCleaner(connector)
	return cl.CleanTopics(connOpts, context.Background(), "test_producer")
}
