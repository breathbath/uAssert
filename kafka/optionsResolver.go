package kafka

import (
	"github.com/breathbath/uAssert/options"
	"time"
)

func ResolveTopicOption(opts options.Options) (string, error) {
	return options.ResolveValueStr("topic", opts, "", true)
}

func ResolveTopicCreateOptions(opts options.Options) (topic string, numPartitions int, replFactor int, timeout time.Duration, err error) {
	topic, err = ResolveTopicOption(opts)
	if err != nil {
		return
	}

	numPartitions, err = options.ResolveValueInt("partitions_number", opts, 1, false)
	if err != nil {
		return
	}

	replFactor, err = options.ResolveValueInt("replication_factor", opts, 1, false)
	if err != nil {
		return
	}

	timeoutInt, err := options.ResolveValueInt("timeout_second", opts, 1, false)
	if err != nil {
		return
	}

	timeout = time.Duration(timeoutInt) * time.Second

	return
}

func ResolveConsumerOptions(opts options.Options) (co ConsumerOptions, err error)  {
	co = ConsumerOptions{}

	co.addr, err = ResolveKafkaAddressOptions(opts)
	if err != nil {
		return
	}

	co.readDeadLineSec, err = options.ResolveValueInt("readDeadlineSec", opts, 3, false)

	return
}

func ResolveProducerOptions(opts options.Options) (po ProducerOptions, err error)  {
	po = ProducerOptions{}

	po.addr, err = ResolveKafkaAddressOptions(opts)
	if err != nil {
		return
	}

	po.writeDeadLineSec, err = options.ResolveValueInt("writeDeadLineSec", opts, 3, false)

	return
}


func ResolveKafkaAddressOptions(opts options.Options) (addr Address, err error) {
	var topic string

	topic, err = ResolveTopicOption(opts)
	if err != nil {
		return
	}

	var partition int

	partition, err = options.ResolveValueInt("partition", opts, 0, false)
	if err != nil {
		return
	}

	addr = Address{
		Topic: topic,
		Partition: partition,
	}

	return
}