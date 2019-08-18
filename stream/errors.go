package stream

import (
	"context"
	"github.com/breathbath/go_utils/utils/errs"
	io2 "github.com/breathbath/go_utils/utils/io"
	"io"
)

func CollectErrors(
	errCont *errs.ErrorContainer,
	ctx context.Context,
	errChan chan error,
	withLog bool,
) {
	go func(errCont *errs.ErrorContainer) {
		for {
			select {
			case err := <-errChan:
				if withLog {
					io2.OutputInfo("", "Received err %v", err)
				}
				if err != io.EOF && err != context.Canceled {
					errCont.AddError(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}(errCont)
}
