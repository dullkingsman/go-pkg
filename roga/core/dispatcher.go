package core

import (
	"github.com/google/uuid"
)

type DefaultDispatcher struct{ Dispatcher }

func (d DefaultDispatcher) AddToOperationQueue(operations []Operation, queue *chan<- uuid.UUID) {
	for _, operation := range operations {
		*queue <- operation.Id
	}
}

func (d DefaultDispatcher) AddToLogQueue(logs []Log, queue *chan<- uuid.UUID) {
	for _, log := range logs {
		*queue <- log.Id
	}
}

func (d DefaultDispatcher) DispatchOperations(operations []Operation, writingChannels *WritingChannels) []uuid.UUID {
	if writingChannels == nil {
		var returnable = make([]uuid.UUID, len(operations))

		for i, operation := range operations {
			returnable[i] = operation.Id
		}

		return returnable
	}

	for _, operation := range operations {
		writingChannels.Stdout <- operation
		writingChannels.File <- operation
		writingChannels.External <- operation
	}

	return nil
}

func (d DefaultDispatcher) DispatchLogs(logs []Log, channels *WritingChannels) []uuid.UUID {
	if channels == nil {
		var returnable = make([]uuid.UUID, len(logs))

		for i, log := range logs {
			returnable[i] = log.Id
		}

		return returnable
	}

	for _, log := range logs {
		channels.Stdout <- log
		channels.File <- log
		channels.External <- log
	}

	return nil
}
