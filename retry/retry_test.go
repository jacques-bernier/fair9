package retry

import (
	"context"
	"errors"
	"testing"
)

func TestExecute_DefaultNoRetries(t *testing.T) {
	retrier := NewRetry()

	expectedError := errors.New("oh no!")

	callCount := 0
	err := retrier.Execute(context.Background(), func(c context.Context) error {
		callCount += 1
		return expectedError
	})

	if callCount > 1 {
		t.Fatal("expected a single call to the function")
	}

	if !errors.Is(err, expectedError) {
		t.Fatalf("expected %s to be returbed", expectedError)
	}
}

func TestExecute_Tokens(t *testing.T) {
	retrier := NewRetry(WithRefillEvery(10), WithMaxTokenBalance(10))

	// No tokens first
	callCount := 0
	err := retrier.Execute(context.Background(), func(c context.Context) error {
		callCount += 1
		return errors.New("oh no!")
	})

	if callCount > 1 {
		t.Fatal("expected a single call to the function since there are no tokens initially")
	}

	if err == nil {
		t.Fatal("expected an error")
	}

	callCount = 0
	runs := 101
	for i := 0; i < runs; i++ {
		err := retrier.Execute(context.Background(), func(c context.Context) error {
			callCount += 1
			return nil
		})
		if err != nil {
			t.Fatal("expected no errors")
		}
	}

	if callCount != runs {
		t.Fatalf("expected %d got %d calls", runs, callCount)
	}

	// Now we have some tokens..
	callCount = 0
	runs = 100
	for i := 0; i < runs; i++ {
		err := retrier.Execute(context.Background(), func(c context.Context) error {
			callCount += 1
			return errors.New("oh no!")
		})
		if err == nil {
			t.Fatal("expected an error")
		}
	}

	// only 10 of the calls should be retried since we have a max of 10 retry tokens
	expectedCalls := int64(runs) + 10
	if int64(callCount) != expectedCalls {
		t.Fatalf("expected exactly %d calls, got %d", expectedCalls, callCount)
	}
}

func TestExecute_ContextCancellation(t *testing.T) {
	// This test ensures that if the context is cancelled, we don't run retry
	retrier := NewRetry(WithInitialRetryTokenBalance(1))

	callCount := 0
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	retrier.Execute(ctx, func(c context.Context) error {
		callCount += 1
		return errors.New("err")
	})

	if callCount > 1 {
		t.Fatal("expected no retries since the context is cancelled")
	}
}
