package simulation

import (
	"context"
	"github.com/breathbath/uAssert/options"
)

type KafkaCleaner struct {
	connector *KafkaConnector
}

func NewKafkaCleaner(connector *KafkaConnector) *KafkaCleaner{
	return &KafkaCleaner{connector: connector}
}

func (kp *KafkaCleaner) CleanTopics(opts options.Options, con context.Context, topics ...string) error {
	kafkaOpts, err := ResolveKafkaOptions(opts)
	if err != nil {
		return err
	}

	conn, err := kp.connector.GetConn(kafkaOpts, con)
	if err != nil {
		return err
	}

	return conn.DeleteTopics(topics...)
}


