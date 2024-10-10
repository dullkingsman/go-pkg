package utils

import (
	"context"
	"sync"
)

// SliceContains checks if a slice of any comparable type contains the specified value.
func SliceContains[T comparable](array []T, value T, compareFunc ...func(T, T) bool) bool {
	var compare = func(a, b T) bool {
		return a == b
	}

	if len(compareFunc) > 0 {
		compare = compareFunc[0]
	}

	for _, element := range array {
		if compare(element, value) {
			return true
		}
	}

	return false
}

// ConcurrentIterSlice iterates over a slice concurrently and calls the specified
// callback for each element.
//
// ___
//
// NOTE: It uses the default context to cancel the iteration.
func ConcurrentIterSlice[T any](slice []T) func(func(int, T) bool) {
	return ConcurrentIterSliceContext(context.Background(), slice)
}

// ConcurrentIterSliceContext iterates over a slice concurrently and calls the specified
// callback for each element.
//
// ___
//
// NOTE: It uses the specified context to cancel the iteration.
func ConcurrentIterSliceContext[T any](ctx context.Context, slice []T) func(func(int, T) bool) {
	return func(yield func(int, T) bool) {
		var ctx, cancel = context.WithCancel(ctx)

		defer cancel()

		var wg sync.WaitGroup

		wg.Add(len(slice))

		for index, element := range slice {
			go func() {
				defer wg.Done()

				select {
				case <-ctx.Done():
					return
				default:
					if !yield(index, element) {
						cancel()
					}
				}
			}()
		}

		wg.Wait()
	}
}
