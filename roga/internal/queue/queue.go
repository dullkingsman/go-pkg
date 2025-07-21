package queue

import (
	"log"
	"sync"
)

type SelfConsumingQueue[T any] struct {
	ch           chan T
	capacity     int
	stopCh       chan struct{}
	flushCh      chan bool
	mu           sync.RWMutex
	stopped      bool
	dependencies []*SelfConsumingQueue[T]
	dependents   []*SelfConsumingQueue[T]
	wg           *sync.WaitGroup
}

func NewQueue[T any](
	fn func([]T),
	capacity int,
	dependencies []*SelfConsumingQueue[T],
	dependents []*SelfConsumingQueue[T],
	wg *sync.WaitGroup,
) *SelfConsumingQueue[T] {
	var q = &SelfConsumingQueue[T]{
		ch:           make(chan T, capacity),
		capacity:     capacity,
		stopCh:       make(chan struct{}),
		flushCh:      make(chan bool),
		dependencies: dependencies,
		dependents:   dependents,
	}

	if wg != nil {
		wg.Add(1)
	} else {
		q.wg = &sync.WaitGroup{}
		q.wg.Add(1)
	}

	go func() {
		if wg != nil {
			defer wg.Done()
		} else {
			defer q.wg.Done()
		}

		q.consume(fn)
	}()

	return q
}

func (q *SelfConsumingQueue[T]) Enqueue(value T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.stopped {
		q.ch <- value
	}
}

func (q *SelfConsumingQueue[T]) EnqueueMany(value []T) {
	if !q.stopped {
		for _, v := range value {
			q.ch <- v
		}
	}
}

func (q *SelfConsumingQueue[T]) Stop() {
	if q.stopped {
		return
	}

	q.flushCh <- true
	close(q.stopCh)
	//close(q.ch)

	q.stopped = true

	if q.wg != nil {
		q.wg.Wait()
	}
}

func (q *SelfConsumingQueue[T]) Flush() {
	log.Println("flushing queue")
	q.flushCh <- true
	log.Println("flushed queue")
}

func (q *SelfConsumingQueue[T]) FlushDependents() {
	for _, dependents := range q.dependents {
		dependents.Flush()
	}
}

func (q *SelfConsumingQueue[T]) StopDependents() {
	for _, dependent := range q.dependents {
		dependent.Stop()
	}
}

func (q *SelfConsumingQueue[T]) Stopped() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.stopped
}

func (q *SelfConsumingQueue[T]) DependenciesAreEmpty() bool {
	for _, dependency := range q.dependencies {
		if dependency.Size() > 0 {
			return false
		}
	}

	return true
}

func (q *SelfConsumingQueue[T]) Size() int {
	return len(q.ch)
}

func (q *SelfConsumingQueue[T]) Cap() int {
	return q.capacity
}

func (q *SelfConsumingQueue[T]) AddDependencies(dependencies ...*SelfConsumingQueue[T]) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.dependencies = append(q.dependencies, dependencies...)
}

func (q *SelfConsumingQueue[T]) AddDependents(dependents ...*SelfConsumingQueue[T]) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.dependents = append(q.dependents, dependents...)
}

func (q *SelfConsumingQueue[T]) consume(fn func([]T)) {
	defer func() {
		log.Println("recovered from panic in queue consumer")
		return
	}()
	for {
		if q.stopped {
			if breakout := q.handleStopped(fn); breakout {
				return
			}
		}

		select {
		case item := <-q.ch:
			fn([]T{item})
		case <-q.stopCh:
			q.mu.Lock()
			q.stopped = true
			q.mu.Unlock()

			q.drainAndHandle(fn)
		case <-q.flushCh:
			if stopped := q.drainAndHandle(fn); stopped {
				var breakout = q.handleStopped(fn)

				if breakout {
					return
				}

				q.mu.Lock()
				q.stopped = true
				q.mu.Unlock()
			}

			q.FlushDependents()
		}
	}
}

func (q *SelfConsumingQueue[T]) drainAndHandle(fn func([]T)) bool {
	var drained, stopped = q.drain()

	if len(drained) > 0 {
		fn(drained)
	}

	return stopped
}

func (q *SelfConsumingQueue[T]) handleStopped(fn func([]T)) bool {
	q.drainAndHandle(fn)

	if !q.DependenciesAreEmpty() {
		return false
	}

	q.StopDependents()

	return true
}

func (q *SelfConsumingQueue[T]) drain() ([]T, bool) {
	var items []T

	for {
		select {
		case item := <-q.ch:
			items = append(items, item)
		case <-q.stopCh:
			return items, true
		case <-q.flushCh:
			return items, false
		default:
			return items, false
		}
	}
}
