package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/breathbath/go_utils/utils/errs"
	"github.com/breathbath/go_utils/utils/io"
	"github.com/breathbath/uAssert/options"
	"sync"
	"time"
)

type StreamFacade struct {
	connStr       string
	clusterAdmin  sarama.ClusterAdmin
	publisher     sarama.SyncProducer
	consumer      sarama.Consumer
	partConsumers sync.Map
}

func NewStreamFacade(connStr string) (*StreamFacade, error) {
	conf := sarama.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	conf.Version = sarama.V2_0_0_0

	clusterAdmin, err := sarama.NewClusterAdmin([]string{connStr}, conf)

	return &StreamFacade{clusterAdmin: clusterAdmin, connStr: connStr, partConsumers: sync.Map{}}, err
}

func (kf *StreamFacade) PrepareStream(opts options.Options) error {
	topic, numPartitions, replFactor, timeout, err := ResolveTopicCreateOptions(opts)

	if err != nil {
		return err
	}

	return kf.createTopic(topic, numPartitions, replFactor, timeout)
}

func (kf *StreamFacade) Read(opts options.Options, ctx context.Context, errChan chan error, outs chan string) {
	var err error
	consumerOpts, err := ResolveConsumerOptions(opts)
	if err != nil {
		errChan <- err
		return
	}

	if kf.consumer == nil {
		config := sarama.NewConfig()
		config.ClientID = "go-kafka-consumer"
		config.Consumer.Return.Errors = true
		config.Net.ReadTimeout = time.Second * time.Duration(consumerOpts.readDeadLineSec)

		kf.consumer, err = sarama.NewConsumer([]string{kf.connStr}, config)
		if err != nil {
			errChan <- err
			return
		}
	}

	partitionList, err := kf.consumer.Partitions(consumerOpts.topic) //get all partitions on the given topic
	if err != nil {
		errChan <- err
		return
	}

	initialOffset := sarama.OffsetOldest
	for _, partition := range partitionList {
		pc, err := kf.consumer.ConsumePartition(consumerOpts.topic, partition, initialOffset)
		if err != nil {
			errChan <- err
			return
		}

		go func(pc sarama.PartitionConsumer) {
		ConsumerLoop:
			for {
				select {
				case msg := <-pc.Messages():
					io.OutputInfo("", "Got message: '%s'", string(msg.Value))
					outs <- string(msg.Value)
				case cer := <-pc.Errors():
					errChan <- cer
				case <-ctx.Done():
					break ConsumerLoop
				}
			}
		}(pc)
	}

	return
}

func (kf *StreamFacade) PublishMany(opts options.Options, payloads []string) error {
	prodOptions, err := ResolveProducerOptions(opts)
	if err != nil {
		return err
	}

	if kf.publisher == nil {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
		config.ClientID = "go-kafka-consumer"
		config.Producer.Return.Errors = true
		config.Net.WriteTimeout = time.Second * time.Duration(prodOptions.writeDeadLineSec)

		kf.publisher, err = sarama.NewSyncProducer([]string{kf.connStr}, config)
		if err != nil {
			return err
		}
	}

	for _, payload := range payloads {
		msg := &sarama.ProducerMessage{
			Topic:     prodOptions.topic,
			Value:     sarama.StringEncoder(payload),
		}
		part, off, err := kf.publisher.SendMessage(msg)
		if err != nil {
			return err
		}

		io.OutputInfo("", "Published '%s' to %d partition with %d offset", payload, part, off)
	}

	return nil
}

func (kf *StreamFacade) Close(opts options.Options) error {
	errCont := errs.NewErrorContainer()

	if kf.publisher != nil {
		err := kf.publisher.Close()
		if err != nil {
			errCont.AddError(err)
		}
	}

	kf.partConsumers.Range(func(key, value interface{}) bool {
		pc := value.(sarama.PartitionConsumer)
		err := pc.Close()
		if err != nil {
			errCont.AddError(err)
		}
		return true
	})

	topic, err := ResolveTopicOption(opts)
	if err == nil {
		err = kf.deleteTopic(topic)
		if err != nil {
			errCont.AddError(err)
		}
	} else {
		errCont.AddError(err)
	}

	if kf.clusterAdmin != nil {
		err = kf.clusterAdmin.Close()
		if err != nil {
			errCont.AddError(err)
		}
	}

	return errCont.Result(" ")
}

func (kf *StreamFacade) createTopic(topic string, numPartitions int, replFactor int, timeout time.Duration) error {
	topicDetail := &sarama.TopicDetail{}
	topicDetail.NumPartitions = int32(numPartitions)
	topicDetail.ReplicationFactor = int16(replFactor)
	topicDetail.ConfigEntries = make(map[string]*string)

	err := kf.clusterAdmin.CreateTopic(topic, topicDetail, false)
	switch typedErr := err.(type) {
	case *sarama.TopicError:
		if typedErr.Err == sarama.ErrTopicAlreadyExists {
			err = nil
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (kf *StreamFacade) deleteTopic(topic string) error {
	if topic == "" {
		return fmt.Errorf("No topics provided")
	}

	err := kf.clusterAdmin.DeleteTopic(topic)
	if err != nil {
		return err
	}
	io.OutputInfo("", "Deleted topic %s", topic)
	return nil
}
