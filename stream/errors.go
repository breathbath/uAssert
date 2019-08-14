package stream

import (
	"context"
	"github.com/breathbath/go_utils/utils/errs"
	"io"
)

func CollectErrors(
	errs *errs.ErrorContainer,
	ctx context.Context,
	errChans ...chan error,
) {
	go func() {
		fannedInErrsStream := MergeErrorStreams(errChans...)
		for {
			select {
			case err := <-fannedInErrsStream:
				if err != io.EOF && err != context.Canceled {
					errs.AddError(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
