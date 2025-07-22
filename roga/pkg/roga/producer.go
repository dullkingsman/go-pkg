package roga

import (
	"github.com/dullkingsman/go-pkg/roga/pkg/model"
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/google/uuid"
)

type DefaultProducer struct{}

func (d DefaultProducer) ProduceLog(log *Log) *Log {
	return log
}

func (d DefaultProducer) BeginOperation(operation *Operation) *Operation {
	return operation
}

func (d DefaultProducer) EndOperation(operation *Operation) *Operation {
	return operation
}

type Producer interface {
	ProduceLog(log *Log) *Log
	BeginOperation(operation *Operation) *Operation
	EndOperation(operation *Operation) *Operation
}

type LogArgs struct {
	Priority       *Priority              `json:"priority,omitempty"`
	VerbosityClass *VerbosityClass        `json:"verbosityClass,omitempty"`
	Event          *string                `json:"event,omitempty"`
	Outcome        *string                `json:"outcome,omitempty"`
	Message        string                 `json:"message"`
	Actor          *model.Actor           `json:"actor,omitempty"`
	TracingId      *uuid.UUID             `json:"tracingId,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
}

type AuditLogArgs struct {
	LogArgs
}

type EventLogArgs struct {
	LogArgs
}

func (a LogArgs) ToLog() Log {

	log := getLogFromPool()
	log.Message = a.Message
	log.Actor = utils.ValueOr(a.Actor, model.Actor{})
	log.Data = a.Data
	log.Event = a.Event
	log.Outcome = a.Outcome
	log.VerbosityClass = utils.ValueOr(a.VerbosityClass, 0)
	log.Priority = utils.ValueOr(a.Priority, 0)
	//return Log{
	//	Message:        a.Message,
	//	Actor:          utils.ValueOr(a.Actor, model.Actor{}),
	//	Data:           a.Data,
	//	Event:          a.Event,
	//	Outcome:        a.Outcome,
	//	VerbosityClass: utils.ValueOr(a.VerbosityClass, 0),
	//	Priority:       utils.ValueOr(a.Priority, 0),
	//}
	return *log
}

type OperationArgs struct {
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	Actor       *model.Actor `json:"actor"`
}

func (a OperationArgs) ToOperation() Operation {
	return Operation{
		Name:        a.Name,
		Description: a.Description,
		Actor:       utils.ValueOr(a.Actor, model.Actor{}),
	}
}
