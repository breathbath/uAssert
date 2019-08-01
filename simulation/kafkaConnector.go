package simulation

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"strings"
)

type KafkaConnector struct {
	conns map[string] *kafka.Conn
}

func NewKafkaConnector() * KafkaConnector{
	return &KafkaConnector{
		conns: make(map[string] *kafka.Conn),
	}
}

func (kc *KafkaConnector) GetConn(opts KafkaOptions, cont context.Context) (*kafka.Conn, error) {
	conn, ok := kc.conns[opts.topic]
	if ok {
		return conn, nil
	}

	var err error

	kc.conns[opts.topic], err = kafka.DialLeader(
		cont,
		"tcp",
		opts.connStr,
		opts.topic,
		opts.partition,
	)

	return kc.conns[opts.topic], err
}

func (kc *KafkaConnector) Destroy() error {
	errs := []string{}
	for _, conn := range kc.conns {
		err := conn.Close()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("Disconnection failure: %s", strings.Join(errs, ", "))
}
