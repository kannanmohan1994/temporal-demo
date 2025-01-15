package helpers

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

type SimpleHeartbit struct {
	cancelc chan struct{}
}

func StartHeartbit(ctx context.Context, period time.Duration) (h *SimpleHeartbit) {
	h = &SimpleHeartbit{make(chan struct{})}
	go func() {
		for {
			select {
			case <-h.cancelc:
				return
			default:
				{
					activity.RecordHeartbeat(ctx)
					time.Sleep(period)
				}
			}
		}
	}()
	return h
}

func (h *SimpleHeartbit) Stop() {
	close(h.cancelc)
}
