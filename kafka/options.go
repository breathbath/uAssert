package kafka

import (
	"fmt"
)

type ConnOptions struct {
	connStr string
}

func (co ConnOptions) Validate() error {
	if co.connStr == "" {
		return fmt.Errorf("Empty conn string")
	}

	return nil
}

type ProducerOptions struct {
	topic            string
	writeDeadLineSec int
}

func (po ProducerOptions) Validate() error {
	if po.topic == "" {
		return fmt.Errorf("Empty topic")
	}
	return nil
}

type ConsumerOptions struct {
	topic            string
	readDeadLineSec int
}

func (co ConsumerOptions) Validate() error {
	if co.topic == "" {
		return fmt.Errorf("Empty topic")
	}
	return nil
}
