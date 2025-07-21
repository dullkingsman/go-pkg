package roga

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
	"github.com/dullkingsman/go-pkg/roga/internal/queue"
	"github.com/dullkingsman/go-pkg/roga/pkg/model"
	"github.com/dullkingsman/go-pkg/roga/writable"
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/google/uuid"
	"os"
	"time"
)

type Operation struct {
	r                     *Roga
	Id                    uuid.UUID                   `json:"id"`
	Name                  string                      `json:"name"`
	Description           *string                     `json:"description,omitempty"`
	BaseOperationId       *uuid.UUID                  `json:"baseOperationId,omitempty"`
	ParentId              *uuid.UUID                  `json:"parentId,omitempty"`
	EssentialMeasurements model.EssentialMeasurements `json:"essentialMeasurements"`
	Measurements          map[string]interface{}      `json:"measurements,omitempty"`
	Actor                 model.Actor                 `json:"actor"`
	Context               *model.Context              `json:"context,omitempty"`
}

type MeasurementHandler func(*Operation)

func (o Operation) String(formatter writable.Formatter) string {
	return formatter.Format(o)
}

func (o Operation) Json() ([]byte, error) {
	var mj = json.NewManualJson()

	mj.WriteUuidField("id", o.Id)
	mj.WriteStringField("name", o.Name)
	mj.WriteStringField("description", utils.ValueOr(o.Description, ""), true)
	mj.WriteUuidField("baseOperationId", utils.ValueOr(o.BaseOperationId, uuid.Nil), true)
	mj.WriteUuidField("parentId", utils.ValueOr(o.ParentId, uuid.Nil), true)

	var err = mj.WriteJsonField("measurements", o.Measurements, false, true)

	if err != nil {
		return nil, err
	}

	essentialMeasurements, err := o.EssentialMeasurements.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("essentialMeasurements", essentialMeasurements)

	actor, err := o.Actor.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("actor", actor)

	context, err := o.Context.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("context", context, true)

	return mj.End(), nil
}

func (o Operation) Bson() ([]byte, error) {
	return nil, nil
}

func (o *Operation) RegisterInstance(r *Roga) {
	o.r = r
}

func (o *Operation) RegisterInstanceFromParent(p *Operation) {
	o.r = p.r
}

func (o *Operation) LogFatal(args LogArgs) {
	produceLog(
		o.r.producer,
		LevelFatal,
		TypeNormal,
		args,
		o,
		&o.r.currentSystemMetrics,
		1,
		//o.r.channels.operational.production,
		o.r.productionQueue,
	)

	os.Exit(1)
}

func (o *Operation) LogError(args LogArgs) {
	produceLog(
		o.r.producer,
		LevelError,
		TypeNormal,
		args,
		o,
		&o.r.currentSystemMetrics,
		1,
		//o.r.channels.operational.production,
		o.r.productionQueue,
	)
}

func (o *Operation) LogWarn(args LogArgs) {
	produceLog(
		o.r.producer,
		LevelWarn,
		TypeNormal,
		args,
		o,
		&o.r.currentSystemMetrics,
		1,
		//o.r.channels.operational.production,
		o.r.productionQueue,
	)
}

func (o *Operation) LogInfo(args LogArgs) {
	produceLog(
		o.r.producer,
		LevelInfo,
		TypeNormal,
		args,
		o,
		&o.r.currentSystemMetrics,
		1,
		//o.r.channels.operational.production,
		o.r.productionQueue,
	)
}

func (o *Operation) LogDebug(args LogArgs) {
	produceLog(
		o.r.producer,
		LevelDebug,
		TypeNormal,
		args,
		o,
		&o.r.currentSystemMetrics,
		1,
		//o.r.channels.operational.production,
		o.r.productionQueue,
	)
}

func (o *Operation) AuditAction(args AuditLogArgs) {
	produceLog(
		o.r.producer,
		LevelInfo,
		TypeAudit,
		args.LogArgs,
		o,
		&o.r.currentSystemMetrics,
		1,
		//o.r.channels.operational.production,
		o.r.productionQueue,
	)
}

func (o *Operation) CaptureEvent(args EventLogArgs) {
	produceLog(
		o.r.producer,
		LevelInfo,
		TypeEvent,
		args.LogArgs,
		o,
		&o.r.currentSystemMetrics,
		1,
		//o.r.channels.operational.production,
		o.r.productionQueue,
	)
}

func (o *Operation) BeginOperation(args OperationArgs, measurementInitiator ...MeasurementHandler) *Operation {
	return beginOperation(
		o.r.producer,
		args,
		o,
		&o.r.context,
		//o.r.channels.operational.production,
		o.r.productionQueue,
		measurementInitiator...,
	)
}

func (o *Operation) EndOperation(measurementFinalizer ...MeasurementHandler) {
	endOperation(
		o.r.producer,
		o,
		//o.r.channels.operational.production,
		o.r.productionQueue,
		measurementFinalizer...,
	)
}

func beginOperation(
	producer Producer,
	args OperationArgs,
	parent *Operation,
	context *model.Context,
	//ch chan<- writable.Writable,
	q *queue.SelfConsumingQueue[writable.Writable],
	measurementInitiator ...MeasurementHandler,
) *Operation {
	var _measurementInitiator MeasurementHandler = nil

	if len(measurementInitiator) > 0 {
		_measurementInitiator = measurementInitiator[0]
	}

	var operation = args.ToOperation()

	operation.Id = uuid.New()

	operation.EssentialMeasurements = model.EssentialMeasurements{
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

		operation.Context = parent.Context

		operation.RegisterInstanceFromParent(parent)
	} else {
		operation.Context = context
	}

	if _measurementInitiator != nil {
		_measurementInitiator(&operation)
	}

	var _operation = producer.BeginOperation(&operation)

	//ch <- _operation

	q.Enqueue(_operation)

	return _operation
}

func endOperation(
	producer Producer,
	operation *Operation,
	//ch chan<- writable.Writable,
	q *queue.SelfConsumingQueue[writable.Writable],
	measurementFinalizer ...MeasurementHandler,
) {
	var _measurementFinalizer MeasurementHandler = nil

	if len(measurementFinalizer) > 0 {
		_measurementFinalizer = measurementFinalizer[0]
	}

	if operation == nil {
		return
	}

	operation.EssentialMeasurements.EndTime = time.Now().UTC()

	if _measurementFinalizer != nil {
		_measurementFinalizer(operation)
	}

	//ch <- producer.EndOperation(operation)

	q.Enqueue(producer.EndOperation(operation))
}
