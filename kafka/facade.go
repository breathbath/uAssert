package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"io"
	"time"
)

const (
	MIN_BATCH_READ = 10e3 //fetch 10KB min, 1MB max
	MAX_BATCH_READ = 1e6  //fetch 10KB min, 1MB max
)

type MessageResolver func(payload string) error

type Facade struct {
	connector *Connector
	connStr   string
}

func NewFacade(connector *Connector, connStr string) Facade {
	return Facade{connector: connector, connStr: connStr}
}

func (kf Facade) CleanTopics(topics ...string) error {
	if len(topics) == 0 || topics[0] == "" {
		return fmt.Errorf("No topics provided")
	}

	conn, err := kf.connector.GetGeneralConn()
	if err != nil {
		return err
	}

	return conn.DeleteTopics(topics...)
}

func (kc Facade) Read(opts ConsumerOptions, ctx context.Context) (outs chan string, errs chan error) {
	errs = make(chan error, 1000)
	outs = make(chan string, 1000)

	err := opts.Validate()
	if err != nil {
		close(outs)
		errs <- err
		return
	}

	readerConf := kafka.ReaderConfig{
		Brokers:   []string{kc.connStr},
		Topic:     opts.addr.Topic,
		Partition: opts.addr.Partition,
		MinBytes:  MIN_BATCH_READ,
		MaxBytes:  MAX_BATCH_READ,
	}

	if opts.readDeadLineSec > 0 {
		readerConf.MaxWait = time.Duration(opts.readDeadLineSec) * time.Second
	}

	go func() {
		readr := kafka.NewReader(readerConf)
		var msg kafka.Message
		var procErr error

		defer func(){
			closeErr := readr.Close()
			if closeErr != io.EOF {
				procErr = fmt.Errorf("Read error: %v, batch close error: %v", procErr, closeErr)
				errs <- procErr
				close(outs)
			}
		}()
		for {
			msg, procErr = readr.FetchMessage(ctx)
			if procErr != nil {
				close(outs)
				errs <- procErr
				return
			}
			outs <- string(msg.Value)
		}
	}()

	return
}

func (kp Facade) CreateTopics(con context.Context, topics ...string) error {
	if len(topics) == 0 || topics[0] == "" {
		return fmt.Errorf("No topics provided")
	}

	conn, err := kp.connector.GetGeneralConn()
	if err != nil {
		return err
	}

	topicConfigs := make([]kafka.TopicConfig, 0, len(topics))
	for _, topic := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}

	return nil
}

func (kp Facade) Publish(payload string, opts ProducerOptions, con context.Context) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	conn, err := kp.connector.GetConn(opts.addr, con)
	if err != nil {
		return err
	}

	if opts.writeDeadLineSec > 0 {
		err = conn.SetWriteDeadline(time.Now().Add(time.Duration(opts.writeDeadLineSec) * time.Second))
		if err != nil {
			return err
		}
	}

	_, err = conn.Write([]byte(payload))
	return err
}
