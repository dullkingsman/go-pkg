package writer

import (
	"github.com/dullkingsman/go-pkg/roga/internal/queue"
	"github.com/dullkingsman/go-pkg/roga/pkg/roga"
	"github.com/dullkingsman/go-pkg/roga/writable"
	"github.com/dullkingsman/go-pkg/utils"
	"sync"
)

type QueuedWriter[T writable.Writable] struct {
	writer *SyncedWriter
	q      *queue.SelfConsumingQueue[T]
}

func NewQueuedWriter[T writable.Writable](
	produceWriter func(currentWriter *SyncedWriter) *SyncedWriter,
	formatter writable.Formatter,
	capacity int,
	dependencies []*queue.SelfConsumingQueue[T],
	wg *sync.WaitGroup,
) *QueuedWriter[T] {
	var w = &QueuedWriter[T]{
		writer: produceWriter(nil),
	}

	var handler = func(items []T) {
		var writer = produceWriter(w.writer)

		if writer == nil {
			writer = w.writer
		}

		for _, item := range items {
			if utils.UnderlyingValueIsNil(item) {
				continue
			}

			var value = item.String(formatter) + "\n"
			_, err := writer.WriteString(value)
			//putLogToPool(log)
			if err != nil {
				// TODO what to do
				continue
			}
			if log, ok := any(item).(*roga.Log); ok {
				//putLogToPool(log)
				//
				roga.PutLogFromPool(log)
			}
		}

		writer.Flush() // maybe not here flush up
	}

	w.q = queue.NewQueue(handler, capacity, dependencies, nil, wg)

	return w
}

func (w *QueuedWriter[T]) GetQueue() *queue.SelfConsumingQueue[T] {
	return w.q
}
