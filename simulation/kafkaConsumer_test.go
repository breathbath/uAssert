package simulation

import (
	"fmt"
	"github.com/breathbath/go_utils/utils/env"
	"github.com/breathbath/uAssert/options"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKafkaIsAvailableForConsumer(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	topic := "test_consumer"
	consumer := NewKafkaConsumer(func(msg string) error {
		fmt.Println(msg)
		return nil
	})

	readOptions := options.Options{
		"topic":               topic,
		"conn":                env.ReadEnvOrFail("KAFKA_CONN_STR"),
		"read_dead_line_sec": 0,
	}

	err := consumer.Read(readOptions)
	assert.NoError(t, err)
}
