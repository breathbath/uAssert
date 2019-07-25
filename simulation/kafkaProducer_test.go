package simulation

import (
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
	producer := NewKafkaProducer()
	pubOptions := options.Options{
		"topic":               topic,
		"conn":                env.ReadEnvOrFail("KAFKA_CONN_STR"),
		"write_dead_line_sec": 10,
	}

	err := producer.Publish("some payload", pubOptions)
	assert.NoError(t, err)

	err = producer.CleanTopics("test_producer")
	assert.NoError(t, err)

	err = producer.Destroy()
	assert.NoError(t, err)
}
