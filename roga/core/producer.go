package core

import (
	"github.com/google/uuid"
	"time"
)

type DefaultProducer struct{ Producer }

func (d DefaultProducer) LogFatal(
	args LogArgs,
	operation *Operation,
	context Context,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelFatal,
		args,
		operation,
		context,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) LogError(
	args LogArgs,
	operation *Operation,
	context Context,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelError,
		args,
		operation,
		context,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) LogWarn(
	args LogArgs,
	operation *Operation,
	context Context,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelWarn,
		args,
		operation,
		context,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) LogInfo(
	args LogArgs,
	operation *Operation,
	context Context,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelInfo,
		args,
		operation,
		context,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) LogDebug(
	args LogArgs,
	operation *Operation,
	context Context,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelDebug,
		args,
		operation,
		context,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) BeginOperation(
	args OperationArgs,
	parent *Operation,
	measurementInitiator *MeasurementHandler,
	ch *chan Writable,
) *Operation {
	if ch == nil {
		return nil
	}

	var operation = args.ToOperation()

	operation.Id = uuid.New()

	operation.EssentialMeasurements = EssentialMeasurements{
		StartTime: time.Now().UTC(),
	}

	if parent != nil {
		operation.ParentId = &parent.Id
		operation.BaseOperationId = parent.BaseOperationId
	}

	if measurementInitiator != nil {
		(*measurementInitiator)(&operation.Measurements)
	}

	*ch <- operation

	return &operation
}

func (d DefaultProducer) EndOperation(
	operation *Operation,
	measurementFinalizer *MeasurementHandler,
	ch *chan Writable,
) {
	if operation == nil {
		return
	}

	operation.EssentialMeasurements.EndTime = time.Now().UTC()

	if measurementFinalizer != nil {
		(*measurementFinalizer)(&operation.Measurements)
	}

	*ch <- *operation
}

func produceLog(
	logLevel Level,
	logArgs LogArgs,
	operation *Operation,
	context Context,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	if ch == nil {
		return nil
	}

	var log = logArgs.ToLog()

	log.Level = logLevel

	log.Context = context

	log.Id = uuid.New()

	log.Timestamp = time.Now().UTC()

	log.Stack.Frames = getStackFrames(framesToSkip + 1)

	if operation != nil {
		log.OperationId = operation.Id

		if operation.BaseOperationId != nil {
			log.TracingId = *operation.BaseOperationId
		}
	}

	*ch <- log

	return &log
}
