package simulation

import (
	"context"
	"github.com/breathbath/uAssert/options"
	"github.com/segmentio/kafka-go"
	"time"
)

type PublishingOptions struct {
	topic     string
	connStr   string
	partition int
	writeDeadLineSec int
}

type KafkaProducer struct {
	conn *kafka.Conn
}

func NewKafkaProducer() *KafkaProducer{
	return &KafkaProducer{}
}

func (kp *KafkaProducer) Publish(payload string, opts options.Options) error {
	pubOptions, err := ResolveKafkaPublishingOptions(opts)
	if err != nil {
		return err
	}

	err = kp.initConn(pubOptions)
	if err != nil {
		return err
	}

	if pubOptions.writeDeadLineSec > 0 {
		err = kp.conn.SetWriteDeadline(time.Now().Add(time.Duration(pubOptions.writeDeadLineSec)*time.Second))
		if err != nil {
			return err
		}
	}

	_, err = kp.conn.Write([]byte(payload))
	return err
}

func (kp *KafkaProducer) initConn(pubOptions PublishingOptions) error {
	var err error

	kp.conn, err = kafka.DialLeader(
		context.Background(),
		"tcp",
		pubOptions.connStr,
		pubOptions.topic,
		pubOptions.partition,
	)

	return err
}

func (kp *KafkaProducer) Destroy() error {
	if kp.conn != nil {
		return kp.conn.Close()
	}

	return nil
}

func (kp *KafkaProducer) CleanTopics(topics ...string) error {
	if kp.conn == nil {
		return nil
	}

	return kp.conn.DeleteTopics(topics...)
}


