package retry

import (
	"context"
	"errors"
	"fmt"
)

func ExampleFoo() {
	retrier, _ := NewRetry()

	ctx := context.TODO()

	// We don't have retry tokens to start
	callCount := 0
	_ = retrier.Execute(ctx, func(ctx context.Context) error {
		callCount += 1
		return errors.New("oh no!")
	})
	fmt.Println(callCount)

	// Accumulate retry tokens by successfuly completing some operations
	for i := 0; i < 100; i++ {
		_ = retrier.Execute(ctx, func(ctx context.Context) error {
			// do something that does not error
			return nil
		})
	}

	// use your retry tokens
	callCount = 0
	_ = retrier.Execute(ctx, func(ctx context.Context) error {
		callCount += 1
		return errors.New("oh no!")
	})

	fmt.Println(callCount)
	// Output: 1
	// 2
}
