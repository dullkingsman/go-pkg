package roga

import (
	"github.com/google/uuid"
	"time"
)

type DefaultProducer struct{ Producer }

func (d DefaultProducer) LogFatal(
	args LogArgs,
	operation *Operation,
	currentSystemMetrics SystemMetrics,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelFatal,
		TypeNormal,
		args,
		operation,
		&currentSystemMetrics,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) LogError(
	args LogArgs,
	operation *Operation,
	currentSystemMetrics SystemMetrics,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelError,
		TypeNormal,
		args,
		operation,
		&currentSystemMetrics,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) LogWarn(
	args LogArgs,
	operation *Operation,
	currentSystemMetrics SystemMetrics,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelWarn,
		TypeNormal,
		args,
		operation,
		&currentSystemMetrics,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) LogInfo(
	args LogArgs,
	operation *Operation,
	currentSystemMetrics SystemMetrics,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelInfo,
		TypeNormal,
		args,
		operation,
		&currentSystemMetrics,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) LogDebug(
	args LogArgs,
	operation *Operation,
	currentSystemMetrics SystemMetrics,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelDebug,
		TypeNormal,
		args,
		operation,
		&currentSystemMetrics,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) AuditAction(
	args AuditLogArgs,
	operation *Operation,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelInfo,
		TypeAudit,
		args.LogArgs,
		operation,
		nil,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) CaptureEvent(
	args EventLogArgs,
	operation *Operation,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	return produceLog(
		LevelInfo,
		TypeEvent,
		args.LogArgs,
		operation,
		nil,
		framesToSkip+1,
		ch,
	)
}

func (d DefaultProducer) BeginOperation(
	args OperationArgs,
	parent *Operation,
	context *Context,
	measurementInitiator MeasurementHandler,
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

		if parent.BaseOperationId == nil {
			operation.BaseOperationId = &parent.Id
		} else {
			operation.BaseOperationId = parent.BaseOperationId
		}

		if args.Actor == nil {
			operation.Actor = parent.Actor
		}
	} else {
		operation.Context = context
	}

	if measurementInitiator != nil {
		measurementInitiator(&operation.Measurements)
	}

	*ch <- operation

	return &operation
}

func (d DefaultProducer) EndOperation(
	operation *Operation,
	measurementFinalizer MeasurementHandler,
	ch *chan Writable,
) {
	if operation == nil {
		return
	}

	operation.EssentialMeasurements.EndTime = time.Now().UTC()

	if measurementFinalizer != nil {
		measurementFinalizer(&operation.Measurements)
	}

	*ch <- *operation
}

func produceLog(
	logLevel Level,
	logType Type,
	logArgs LogArgs,
	operation *Operation,
	currentSystemMetrics *SystemMetrics,
	framesToSkip int,
	ch *chan Writable,
) *Log {
	if ch == nil {
		return nil
	}

	var log = logArgs.ToLog()

	log.Level = logLevel
	log.Type = logType

	if currentSystemMetrics != nil {
		log.SystemMetrics = *currentSystemMetrics
	}

	log.Id = uuid.New()

	log.Timestamp = time.Now().UTC()

	log.Stack.Frames = getStackFrames(framesToSkip + 1)

	if operation != nil {
		log.OperationId = operation.Id

		if operation.BaseOperationId != nil {
			log.TracingId = *operation.BaseOperationId
		}

		if logArgs.Actor == nil {
			log.Actor = operation.Actor
		}
	}

	*ch <- log

	return &log
}
