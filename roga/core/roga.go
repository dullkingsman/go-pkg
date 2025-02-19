package core

import (
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/google/uuid"
	"sync"
	"time"
)

func Init(config ...Config) Roga {
	var _config = defaultRogaConfig

	if len(config) > 0 {
		var __config = config[0]

		if __config.Instance != nil {
			_config.Instance = __config.Instance
		}

		if __config.Producer != nil {
			_config.Producer = __config.Producer
		}

		if __config.Monitor != nil {
			_config.Monitor = __config.Monitor
		}

		if __config.Dispatcher != nil {
			_config.Dispatcher = __config.Dispatcher
		}

		if __config.Writer != nil {
			_config.Writer = __config.Writer
		}
	}

	var instance = Roga{
		context:     defaultOperationContext,
		config:      *_config.Instance,
		metricsLock: &sync.RWMutex{},
		producer:    _config.Producer,
		monitor:     _config.Monitor,
		dispatcher:  _config.Dispatcher,
		writer:      _config.Writer,
		rootOperation: Operation{
			Id:          uuid.New(),
			Name:        "root",
			Description: utils.PtrOf("A program run!"),
			EssentialMeasurements: EssentialMeasurements{
				StartTime: time.Now().UTC(),
			},
			Actor: Actor{Type: 1},
		},
		logQueue:          make(chan uuid.UUID, _config.Instance.maxLogQueueSize),
		operationQueue:    make(chan uuid.UUID, _config.Instance.maxOperationQueueSize),
		productionChannel: make(chan Writable),
		writingChannels: WritingChannels{
			Stdout:   make(chan Writable),
			File:     make(chan Writable),
			External: make(chan Writable),
		},
	}

	instance.context.Environment.ApplicationEnvironment.ServiceName = _config.ServiceName

	instance.rootOperation.r = &instance

	return instance
}

func (r *Roga) LogFatal(args LogArgs) {
	r.producer.LogFatal(
		args,
		&r.rootOperation,
		r.context,
		1,
		&r.productionChannel,
	)
}

func (r *Roga) LogError(args LogArgs) {
	r.producer.LogError(
		args,
		&r.rootOperation,
		r.context,
		1,
		&r.productionChannel,
	)
}

func (r *Roga) LogWarn(args LogArgs) {
	r.producer.LogWarn(
		args,
		&r.rootOperation,
		r.context,
		1,
		&r.productionChannel,
	)
}

func (r *Roga) LogInfo(args LogArgs) {
	r.producer.LogInfo(
		args,
		&r.rootOperation,
		r.context,
		1,
		&r.productionChannel,
	)
}

func (r *Roga) LogDebug(args LogArgs) {
	r.producer.LogDebug(
		args,
		&r.rootOperation,
		r.context,
		1,
		&r.productionChannel,
	)
}

func (r *Roga) BeginOperation(args OperationArgs, measurementInitiator ...MeasurementHandler) *Operation {
	var _measurementInitiator MeasurementHandler = nil

	if len(measurementInitiator) > 0 {
		_measurementInitiator = measurementInitiator[0]
	}

	var operation = r.producer.BeginOperation(
		args,
		&r.rootOperation,
		&_measurementInitiator,
		&r.productionChannel,
	)

	operation.r = r

	return operation
}

func (o *Operation) LogFatal(args LogArgs) {
	var log = o.r.producer.LogFatal(
		args,
		o,
		o.r.context,
		1,
		&o.r.productionChannel,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogError(args LogArgs) {
	var log = o.r.producer.LogError(
		args,
		o,
		o.r.context,
		1,
		&o.r.productionChannel,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogWarn(args LogArgs) {
	var log = o.r.producer.LogWarn(
		args,
		o,
		o.r.context,
		1,
		&o.r.productionChannel,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogInfo(args LogArgs) {
	var log = o.r.producer.LogInfo(
		args,
		o,
		o.r.context,
		1,
		&o.r.productionChannel,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogDebug(args LogArgs) {
	var log = o.r.producer.LogDebug(
		args,
		o,
		o.r.context,
		1,
		&o.r.productionChannel,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) BeginOperation(args OperationArgs, measurementInitiator ...MeasurementHandler) *Operation {
	var _measurementInitiator MeasurementHandler = nil

	if len(measurementInitiator) > 0 {
		_measurementInitiator = measurementInitiator[0]
	}

	var operation = o.r.producer.BeginOperation(
		args,
		o,
		&_measurementInitiator,
		&o.r.productionChannel,
	)

	o.OperationChildren = append(o.OperationChildren, operation.Id)

	operation.r = o.r

	return operation
}

func (o *Operation) EndOperation(measurementFinalizer ...MeasurementHandler) {
	var _measurementFinalizer MeasurementHandler = nil

	if len(measurementFinalizer) > 0 {
		_measurementFinalizer = measurementFinalizer[0]
	}

	o.r.producer.EndOperation(
		o,
		&_measurementFinalizer,
		&o.r.productionChannel,
	)
}

func (r *Roga) consumeProductionChannel() {
	for {
		if len(r.productionChannel) < r.config.maxProductionChannelItems {
			continue
		}

		for i := 0; i < r.config.maxProductionChannelItems; i++ {
			var writable = <-r.productionChannel

			if operation, ok := writable.(Operation); ok {
				var _, alreadyInBuffer = r.operationBuffer[operation.Id]

				if !alreadyInBuffer {
					continue
				}

				r.operationBuffer[operation.Id] = operation
				r.operationQueue <- operation.Id
			} else if log, ok := writable.(Log); ok {
				r.logBuffer[log.Id] = log
				r.logQueue <- log.Id
			}
		}
	}
}

func (r *Roga) consumeOperationQueue() {
	for {
		if len(r.operationQueue) < r.config.maxOperationQueueSize {
			continue
		}

		var operations = make([]Operation, r.config.maxOperationQueueSize)

		for i := 0; i < r.config.maxOperationQueueSize; i++ {

			var operationId = <-r.operationQueue

			var operation, ok = r.operationBuffer[operationId]

			if !ok {
				continue
			}

			operations[i] = operation
		}

		r.dispatcher.DispatchOperations(operations, &r.writingChannels)
	}
}

func (r *Roga) consumeLogQueue() {
	for {
		if len(r.logQueue) < r.config.maxLogQueueSize {
			continue
		}

		var logs = make([]Log, r.config.maxLogQueueSize)

		for i := 0; i < r.config.maxLogQueueSize; i++ {

			var logId = <-r.logQueue

			var log, ok = r.logBuffer[logId]

			if !ok {
				continue
			}

			logs[i] = log
		}

		r.dispatcher.DispatchLogs(logs, &r.writingChannels)
	}
}

func (r *Roga) consumeStdoutWrites(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if len(r.writingChannels.Stdout) < r.config.maxStdoutWriterChannelItems {
			continue
		}

		var operations []Operation

		var logs []Log

		for i := 0; i < r.config.maxStdoutWriterChannelItems; i++ {
			var writable = <-r.writingChannels.Stdout

			if operation, ok := writable.(Operation); ok {
				operations = append(operations, operation)
			} else if log, ok := writable.(Log); ok {
				logs = append(logs, log)
			}

		}

		if len(operations) > 0 {
			r.writer.WriteOperationsToStdout(operations, r)
		}

		if len(logs) > 0 {
			r.writer.WriteLogsToStdout(logs, r)
		}
	}
}

func (r *Roga) consumeFileWrites(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if len(r.writingChannels.File) < r.config.maxFileWriterChannelItems {
			continue
		}

		var operations []Operation

		var logs []Log

		for i := 0; i < r.config.maxFileWriterChannelItems; i++ {
			var writable = <-r.writingChannels.File

			if operation, ok := writable.(Operation); ok {
				operations = append(operations, operation)
			} else if log, ok := writable.(Log); ok {
				logs = append(logs, log)
			}
		}

		if len(operations) > 0 {
			r.writer.WriteOperationsToFile(operations, r)
		}

		if len(logs) > 0 {
			r.writer.WriteLogsToFile(logs, r)
		}
	}
}

func (r *Roga) consumeExternalWrites(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if len(r.writingChannels.External) < r.config.maxExternalWriterChannelItems {
			continue
		}

		var operations []Operation

		var logs []Log

		for i := 0; i < r.config.maxExternalWriterChannelItems; i++ {
			var writable = <-r.writingChannels.External

			if operation, ok := writable.(Operation); ok {
				operations = append(operations, operation)
			} else if log, ok := writable.(Log); ok {
				logs = append(logs, log)
			}
		}

		if len(operations) > 0 {
			r.writer.WriteOperationsToExternal(operations, r)
		}

		if len(logs) > 0 {
			r.writer.WriteLogsToExternal(logs, r)
		}
	}
}

func (r *Roga) consumeChannels() {
	go r.monitorAndUpdateSystemMetrics()

	go r.consumeProductionChannel()

	go r.consumeLogQueue()

	go r.consumeOperationQueue()

	go func() {
		var wg sync.WaitGroup

		for i := 0; i < r.config.maxStdoutWriters; i++ {
			go r.consumeStdoutWrites(&wg)
		}

		wg.Wait()
	}()

	go func() {
		var wg sync.WaitGroup

		for i := 0; i < r.config.maxFileWriters; i++ {
			go r.consumeFileWrites(&wg)
		}

		wg.Wait()
	}()

	go func() {
		var wg sync.WaitGroup

		for i := 0; i < r.config.maxExternalWriters; i++ {
			go r.consumeExternalWrites(&wg)
		}

		wg.Wait()
	}()
}
