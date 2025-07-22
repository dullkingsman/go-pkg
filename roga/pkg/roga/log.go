package roga

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
	"github.com/dullkingsman/go-pkg/roga/internal/queue"
	"github.com/dullkingsman/go-pkg/roga/internal/stacktrace"
	"github.com/dullkingsman/go-pkg/roga/pkg/model"
	"github.com/dullkingsman/go-pkg/roga/writable"
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Log struct {
	Id             uuid.UUID              `json:"id"`
	Type           Type                   `json:"type"`
	Event          *string                `json:"event,omitempty"`
	Outcome        *string                `json:"outcome,omitempty"`
	Level          Level                  `json:"level"`
	Priority       Priority               `json:"priority"`
	VerbosityClass VerbosityClass         `json:"verbosityClass"`
	Message        string                 `json:"message"`
	TracingId      uuid.UUID              `json:"tracingId"`
	OperationId    uuid.UUID              `json:"operationId"`
	Timestamp      time.Time              `json:"timestamp"`
	Stack          model.StackTrace       `json:"stack"`
	Actor          model.Actor            `json:"actor"`
	SystemMetrics  model.SystemMetrics    `json:"systemMetrics"`
	Data           map[string]interface{} `json:"data,omitempty"`
}

var logPool = sync.Pool{
	New: func() interface{} {
		return &Log{
			Data: make(map[string]interface{}),
		}
	},
}

func getLogFromPool() *Log {
	log := logPool.Get().(*Log)
	log.Id = uuid.Nil
	return log
}

func PutLogFromPool(log *Log) {
	logPool.Put(log)
}

type (
	Type           int
	Level          int
	VerbosityClass int
	Priority       int
)

const (
	TypeNormal Type = 0
	TypeAudit  Type = 1
	TypeEvent  Type = 2

	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
	LevelFatal Level = 12

	VerbosityClassMandatory VerbosityClass = 0
	VerbosityClass1         VerbosityClass = 1
	VerbosityClass2         VerbosityClass = 2
	VerbosityClass3         VerbosityClass = 3
	VerbosityClass4         VerbosityClass = 4
	VerbosityClass5         VerbosityClass = 5

	PriorityOptional Priority = -4
	PriorityLow      Priority = -2
	PriorityMedium   Priority = 0
	PriorityHigh     Priority = 2
	PriorityCritical Priority = 4
)

func (l Log) String(formatter writable.Formatter) string {
	return formatter.Format(l)
}

func (l Log) Json() ([]byte, error) {
	var mj = json.NewManualJson()

	mj.WriteUuidField("id", l.Id)
	mj.WriteInt64Field("type", int64(l.Type))

	mj.WriteStringField("event", utils.ValueOr(l.Event, ""), true)
	mj.WriteStringField("outcome", utils.ValueOr(l.Outcome, ""), true)

	mj.WriteInt64Field("level", int64(l.Level))
	mj.WriteInt64Field("priority", int64(l.Priority))
	mj.WriteInt64Field("verbosityClass", int64(l.VerbosityClass))

	mj.WriteStringField("message", l.Message)

	mj.WriteUuidField("tracingId", l.TracingId, true)
	mj.WriteUuidField("operationId", l.OperationId, true)

	mj.WriteTimeField("timestamp", l.Timestamp, nil, true)

	var stack, err = l.Stack.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("stack", stack)

	actor, err := l.Actor.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("actor", actor)

	systemMetrics, err := l.SystemMetrics.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("systemMetrics", systemMetrics)

	err = mj.WriteJsonField("data", l.Data, false, true)

	if err != nil {
		return nil, err
	}

	return mj.End(), nil
}

func (l Log) Bson() ([]byte, error) {
	return nil, nil
}

func produceLog(
	producer Producer,
	logLevel Level,
	logType Type,
	logArgs LogArgs,
	operation *Operation,
	currentSystemMetrics *model.SystemMetrics,
	framesToSkip int,
//ch chan<- writable.Writable,
	q *queue.SelfConsumingQueue[writable.Writable],
) *Log {
	var log = logArgs.ToLog()

	log.Level = logLevel
	log.Type = logType

	if currentSystemMetrics != nil {
		log.SystemMetrics = *currentSystemMetrics
	}

	log.Id = uuid.New()

	log.Timestamp = time.Now().UTC()

	log.Stack.Frames = stacktrace.GetStackFrames(framesToSkip + 1)

	if operation != nil {
		log.OperationId = operation.Id

		if operation.BaseOperationId != nil {
			log.TracingId = *operation.BaseOperationId
		}

		if logArgs.Actor == nil {
			log.Actor = operation.Actor
		}
	}

	var _log = producer.ProduceLog(&log)

	//ch <- _log
	q.Enqueue(_log)

	return _log
}
