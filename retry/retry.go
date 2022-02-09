package retry

import (
	"context"
	"sync/atomic"
)

func NewRetry() (*retrier, error) {
	return &retrier{
		every:    99,
		counter:  0,
		token:    0,
		maxToken: 10,
	}, nil
}

type retrier struct {
	every uint64

	counter  uint64
	token    int64
	maxToken int64
}

func (r *retrier) Execute(ctx context.Context, op func(context.Context) error) error {
	err := op(ctx)
	if err != nil {
		newBalance := atomic.AddInt64(&r.token, -1)
		if newBalance >= 0 {
			return r.Execute(ctx, op)
		}

		// We can't go negative. So we add it back.
		if newBalance == -1 {
			atomic.AddInt64(&r.token, 1)
		}
		return err
	}

	// record success
	newCount := atomic.AddUint64(&r.counter, 1)
	if newCount%r.every == 0 {
		// Add a token
		newToken := atomic.AddInt64(&r.token, 1)

		// Don't let the token balance go above the limit
		if newToken > r.maxToken {
			atomic.AddInt64(&r.token, -1)
		}
	}
	return nil
}
