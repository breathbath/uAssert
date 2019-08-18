package stream

import (
	"context"
	"github.com/breathbath/uAssert/options"
)

type StreamFacade interface {
	PrepareStream(opts options.Options) error
	Read(opts options.Options, ctx context.Context, errs chan error, outs chan string)
	PublishMany(opts options.Options, payloads []string) error
	Close(opts options.Options) error
}
