package simulation

import "github.com/breathbath/uAssert/options"

type KafkaOptions struct {
	topic     string
	connStr   string
	partition int
}

type ProducerKafkaOptions struct {
	KafkaOptions
	writeDeadLineSec int
	maxAttempts int
}

type ConsumerKafkaOptions struct {
	KafkaOptions
	readDeadLineSec int
}

func ResolveKafkaOptions(opts options.Options) (ko KafkaOptions, err error) {
	ko = KafkaOptions{}

	ko.topic, err = options.ResolveValueStr("topic", opts, "", true)
	if err != nil {
		return
	}

	ko.connStr, err = options.ResolveValueStr("conn", opts, "", true)
	if err != nil {
		return
	}

	ko.partition, err = options.ResolveValueInt("partition", opts, 0, false)
	if err != nil {
		return
	}

	return ko, nil
}

func ResolveProducerOptions(opts options.Options) (po ProducerKafkaOptions, err error) {
	var kafkaOptions KafkaOptions
	kafkaOptions, err = ResolveKafkaOptions(opts)
	if err != nil {
		return
	}

	po = ProducerKafkaOptions{
		KafkaOptions: kafkaOptions,
	}

	po.writeDeadLineSec, err = options.ResolveValueInt("write_dead_line_sec", opts, 0, false)
	if err != nil {
		return
	}

	po.maxAttempts, err = options.ResolveValueInt("max_attempts", opts, 1, false)
	if err != nil {
		return
	}

	return po, nil
}

func ResolveConsumerOptions(opts options.Options) (co ConsumerKafkaOptions, err error) {
	var kafkaOptions KafkaOptions
	kafkaOptions, err = ResolveKafkaOptions(opts)
	if err != nil {
		return
	}

	co = ConsumerKafkaOptions{
		KafkaOptions: kafkaOptions,
	}

	co.readDeadLineSec, err = options.ResolveValueInt("read_dead_line_sec", opts, 0, false)
	if err != nil {
		return
	}

	return co, nil
}

