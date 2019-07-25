package simulation

import (
	"github.com/breathbath/uAssert/options"
)

func ResolveKafkaPublishingOptions(opts options.Options) (po PublishingOptions, err error) {
	po = PublishingOptions{}

	po.topic, err = options.ResolveValueStr("topic", opts, "", true)
	if err != nil {
		return
	}

	po.connStr, err = options.ResolveValueStr("conn", opts, "", true)
	if err != nil {
		return
	}

	po.partition, err = options.ResolveValueInt("partition", opts, 0, false)
	if err != nil {
		return
	}

	po.writeDeadLineSec, err = options.ResolveValueInt("write_dead_line_sec", opts, 0, false)
	if err != nil {
		return
	}

	return po, nil
}
