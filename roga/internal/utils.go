package internal

import (
	"github.com/dullkingsman/go-pkg/utils"
	"sync"
)

func ConsumeQueue[T any](
	name string,
	queue <-chan T,
	stopChan <-chan bool,
	flushChan <-chan bool,
	dependencies []<-chan T,
	dependentsFlush []chan<- bool,
	dependentsStop []chan<- bool,
	wg *sync.WaitGroup,
	handler func(items []T),
	dontLogLoggerLogs ...bool,
) {
	var log = len(dontLogLoggerLogs) == 0 || dontLogLoggerLogs[0] == false

	defer wg.Done()

	if log {
		utils.LogInfo("roga:startup", "consuming "+name+"...")
	}

	var (
		stopped = false
		flushed = false
	)

	for {
		var items []T

		if stopped {
		FlushBreak:
			for {
				stopped, flushed = checkQueueLifecycle(name, flushChan, stopChan, log)

				if flushed {
					break FlushBreak
				}

				var drained = drainQueue(queue)

				if len(drained) > 0 {
					handler(drained)
				}

				var stop = true

				for _, dependency := range dependencies {
					if len(dependency) > 0 {
						stop = false
						break
					}
				}

				if stop {
					for idx, dependent := range dependentsFlush {
						dependent <- true
						close(dependent)
						dependentsStop[idx] <- true
						close(dependentsStop[idx])
					}

					if log {
						utils.LogInfo("roga:cleanup", "stopped consuming "+name)
					}

					return
				}
			}
		} else if flushed {
			var drained = drainQueue(queue)

			if len(drained) > 0 {
				handler(drained)
			}

			for _, dependent := range dependentsFlush {
				dependent <- true
			}

			flushed = false
		} else {
			stopped, flushed = checkQueueLifecycle(name, flushChan, stopChan, log)

			var drained = drainQueue(queue)

			if len(drained) > 0 {
				handler(drained)
			}
		}

		if items == nil {
			continue
		}
	}
}

func checkQueueLifecycle(queueName string, flush <-chan bool, stop <-chan bool, log bool) (stopped bool, flushed bool) {
	select {
	case <-stop:
		if !stopped {
			stopped = true

			if log {
				utils.LogInfo("roga:ops", "signaled "+queueName+" consumption to stop")
			}
		}
	case <-flush:
		flushed = true

		if log {
			utils.LogInfo("roga:ops", "signaled "+queueName+" to flush")
		}
	default:
	}

	return stopped, flushed
}

func drainQueue[T any](queue <-chan T) []T {
	var size = len(queue)

	if size > 0 {
		var items = make([]T, size)

		for i := 0; i < size; i++ {
			items[i] = <-queue
		}

		return items
	}

	return nil
}
