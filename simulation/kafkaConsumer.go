package simulation

import (
	"context"
	"fmt"
	"github.com/breathbath/uAssert/options"
	"github.com/segmentio/kafka-go"
	"io"
	"time"
)

const (
	MIN_BATCH_READ=10e3 //fetch 10KB min, 1MB max
	MAX_BATCH_READ=1e6 //fetch 10KB min, 1MB max
)

type MessageResolver func(msg string) error

type KafkaConsumer struct {
	msgResolver   MessageResolver
}

func NewKafkaConsumer(msgResolver MessageResolver) *KafkaConsumer {
	return &KafkaConsumer{
		msgResolver:   msgResolver,
	}
}

func (kc *KafkaConsumer) Read(opts options.Options) error {
	readOptions, err := ResolveConsumerOptions(opts)
	if err != nil {
		return err
	}

	config := kafka.ReaderConfig{
		Brokers:   []string{readOptions.connStr},
		Topic:     readOptions.topic,
		MinBytes:  MIN_BATCH_READ,
		MaxBytes:  MAX_BATCH_READ,
	}

	if readOptions.readDeadLineSec > 0 {
		config.MaxWait = time.Duration(readOptions.readDeadLineSec) * time.Second
	}

	r := kafka.NewReader(config)
	var msg kafka.Message
	var procErr error
	for {
		msg, procErr = r.ReadMessage(context.Background())
		if procErr != nil {
			break
		}
		procErr := kc.msgResolver(string(msg.Value))
		if procErr != nil {
			break
		}
	}

	closeErr := r.Close()
	if closeErr != io.EOF {
		procErr = fmt.Errorf("Read error: %v, batch close error: %v", procErr, closeErr)
	}

	return procErr
}
