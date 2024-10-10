package utils

import (
	"context"
	"sync"
)

// ConcurrentIterMap iterates over a map concurrently and calls the specified
// callback for each element.
//
// ___
//
// NOTE: It uses the default context to cancel the iteration.
func ConcurrentIterMap[T comparable, H any](table map[T]H) func(func(T, H) bool) {
	return ConcurrentIterMapContext(context.Background(), table)
}

// ConcurrentIterMapContext iterates over a map concurrently and calls the specified
// callback for each element.
//
// ___
//
// NOTE: It uses the specified context to cancel the iteration.
func ConcurrentIterMapContext[T comparable, H any](ctx context.Context, slice map[T]H) func(func(T, H) bool) {
	return func(yield func(T, H) bool) {
		var ctx, cancel = context.WithCancel(ctx)

		defer cancel()

		var wg sync.WaitGroup

		wg.Add(len(slice))

		for key, value := range slice {
			go func() {
				defer wg.Done()

				select {
				case <-ctx.Done():
					return
				default:
					if !yield(key, value) {
						cancel()
					}
				}
			}()
		}

		wg.Wait()
	}
}
