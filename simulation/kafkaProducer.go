package simulation

import (
	"context"
	"github.com/breathbath/uAssert/options"
	"github.com/segmentio/kafka-go"
	"time"
)

type KafkaProducer struct {
	connector *KafkaConnector
}

func NewKafkaProducer(connector *KafkaConnector) *KafkaProducer{
	return &KafkaProducer{connector: connector}
}

func (kp *KafkaProducer) Publish(payload string, opts options.Options, cont context.Context) error {
	pubOptions, err := ResolveProducerOptions(opts)
	if err != nil {
		return err
	}

	config := kafka.WriterConfig{
		Brokers:   []string{pubOptions.connStr},
		Topic:     pubOptions.topic,
		MaxAttempts: pubOptions.maxAttempts,
	}

	if pubOptions.writeDeadLineSec > 0 {
		config.WriteTimeout = time.Duration(pubOptions.writeDeadLineSec) * time.Second
	}

	conn, err := kp.connector.GetConn(pubOptions.KafkaOptions, cont)
	if err != nil {
		return err
	}

	if pubOptions.writeDeadLineSec > 0 {
		err = conn.SetWriteDeadline(time.Now().Add(time.Duration(pubOptions.writeDeadLineSec)*time.Second))
		if err != nil {
			return err
		}
	}

	_, err = conn.Write([]byte(payload))
	return err
}

