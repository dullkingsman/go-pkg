package utils

import (
	"context"
	"sync"
)

func ConcurrentIterMap[T comparable, H any](table map[T]H) func(func(T, H) bool) {
	return ConcurrentIterMapContext(context.Background(), table)
}

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
