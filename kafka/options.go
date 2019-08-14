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
	addr             Address
	writeDeadLineSec int
}

func (po ProducerOptions) Validate() error {
	err := po.addr.Validate()
	return err
}

type ConsumerOptions struct {
	addr            Address
	readDeadLineSec int
}

func (co ConsumerOptions) Validate() error {
	err := co.addr.Validate()
	return err
}

type Address struct {
	Topic     string
	Partition int
}

func (ka Address) Validate() error {
	if ka.Topic == "" {
		return fmt.Errorf("Empty topic")
	}

	return nil
}
