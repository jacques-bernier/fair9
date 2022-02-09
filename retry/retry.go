package retry

import (
	"context"
	"sync/atomic"
)

type option func(*retrier) *retrier

func NewRetry(options ...option) *retrier {
	r := &retrier{
		every:             0, // No retries
		counter:           0,
		retryTokenBalance: 0,
		maxTokenBalance:   0,
	}

	for _, opt := range options {
		r = opt(r)
	}

	return r
}

func WithMaxTokenBalance(max int64) option {
	return func(r *retrier) *retrier {
		r.maxTokenBalance = max
		return r
	}
}

func WithInitialRetryTokenBalance(balance int64) option {
	return func(r *retrier) *retrier {
		r.retryTokenBalance = balance
		return r
	}
}

// WithRefillEvery sets the parameter that controls how many successful operations
// are required to gain a new retry token.
func WithRefillEvery(every uint64) option {
	return func(r *retrier) *retrier {
		r.every = every
		return r
	}
}

type retrier struct {
	every             uint64
	counter           uint64
	retryTokenBalance int64
	maxTokenBalance   int64
}

// Execute runs the `op` and will use a retry in case of an error
func (r *retrier) Execute(ctx context.Context, op func(context.Context) error) error {
	err := op(ctx)
	// TODO check is the error is retriable using a function provided by the user.
	// TODO add max attempts as a feature to prevent a single call to execute from consuming all tokens
	if err != nil {
		if ctx.Err() != nil {
			return err
		}
		newBalance := atomic.AddInt64(&r.retryTokenBalance, -1)
		if newBalance >= 0 {
			return r.Execute(ctx, op)
		}

		// We can't go negative. So we add it back.
		if newBalance == -1 {
			atomic.AddInt64(&r.retryTokenBalance, 1)
		}
		return err
	}

	// record success
	newCount := atomic.AddUint64(&r.counter, 1)
	if r.every > 0 && newCount%r.every == 0 {
		// Add a token
		newToken := atomic.AddInt64(&r.retryTokenBalance, 1)

		// Don't let the token balance go above the limit
		if newToken > r.maxTokenBalance {
			atomic.AddInt64(&r.retryTokenBalance, -1)
		}
	}
	return nil
}
