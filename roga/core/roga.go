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
		context:         defaultOperationContext,
		config:          *_config.Instance,
		metricsLock:     &sync.RWMutex{},
		consumptionSync: &sync.WaitGroup{},
		started:         false,
		producer:        _config.Producer,
		monitor:         _config.Monitor,
		dispatcher:      _config.Dispatcher,
		writer:          _config.Writer,
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

func (r *Roga) Start() {
	if !r.started {
		r.started = true
		r.consumeChannels()
	}
}

func (r *Roga) Recover() {
	r.writingChannelsFlush.Stdout <- true
	r.writingChannelsFlush.File <- true
	r.writingChannelsFlush.External <- true

	r.operationQueueFlush <- true
	r.logQueueFlush <- true

	r.productionChannelFlush <- true

	r.consumptionSync.Wait()

	utils.LogInfo("roga:signal-handler", "flushed everything")
}

func (r *Roga) Wait() {
	r.consumptionSync.Wait()
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
	r.consumptionSync.Add(1)
	defer r.consumptionSync.Done()

	var flushed = false

	for {
		if flushed {
			for {
				select {
				case writable := <-r.productionChannel:
					addWritableToQueue(writable, r)
				default:
					close(r.productionChannel)

					r.operationQueueFlush <- true
					r.logQueueFlush <- true

					return
				}
			}
		}

		select {
		case <-r.productionChannelFlush:
			if !flushed {
				flushed = true
				continue
			}
		default:
			if len(r.productionChannel) < r.config.maxProductionChannelItems {
				continue
			}

			for i := 0; i < r.config.maxProductionChannelItems; i++ {
				addWritableToQueue(<-r.productionChannel, r)
			}
		}
	}
}

func (r *Roga) consumeOperationQueue() {
	r.consumptionSync.Add(1)
	defer r.consumptionSync.Done()

	var flushed = false

	for {
		var operations []Operation

		if flushed {
			for {
				select {
				case operationId := <-r.operationQueue:
					var operation, ok = r.operationBuffer[operationId]

					if !ok {
						continue
					}

					operations = append(operations, operation)
				default:
					var stop = len(r.productionChannel) == 0

					r.writingChannelsFlush.Stdout <- true
					r.writingChannelsFlush.File <- true
					r.writingChannelsFlush.External <- true

					if stop {
						close(r.operationQueue)
						return
					}

					break
				}
			}
		}

		select {
		case <-r.operationQueueFlush:
			if !flushed {
				flushed = true
				continue
			}
		default:
			if len(r.operationQueue) < r.config.maxOperationQueueSize {
				continue
			}

			operations = make([]Operation, r.config.maxOperationQueueSize)

			for i := 0; i < r.config.maxOperationQueueSize; i++ {
				var operationId = <-r.operationQueue

				var operation, ok = r.operationBuffer[operationId]

				if !ok {
					continue
				}

				operations[i] = operation
			}
		}

		r.dispatcher.DispatchOperations(operations, &r.writingChannels)
		operations = make([]Operation, 0)
	}
}

func (r *Roga) consumeLogQueue() {
	r.consumptionSync.Add(1)
	defer r.consumptionSync.Done()

	var flushed = false

	for {
		var logs []Log

		if flushed {
			for {
				select {
				case logId := <-r.logQueue:
					var log, ok = r.logBuffer[logId]

					if !ok {
						continue
					}

					logs = append(logs, log)
				default:
					var stop = len(r.productionChannel) == 0

					r.writingChannelsFlush.Stdout <- true
					r.writingChannelsFlush.File <- true
					r.writingChannelsFlush.External <- true

					if stop {
						close(r.logQueue)
						return
					}

					break
				}
			}
		}

		select {
		case <-r.logQueueFlush:
			if !flushed {
				flushed = true
				continue
			}
		default:
			if len(r.logQueue) < r.config.maxLogQueueSize {
				continue
			}

			logs = make([]Log, r.config.maxLogQueueSize)

			for i := 0; i < r.config.maxLogQueueSize; i++ {
				var logId = <-r.logQueue

				var log, ok = r.logBuffer[logId]

				if !ok {
					continue
				}

				logs[i] = log
			}
		}

		r.dispatcher.DispatchLogs(logs, &r.writingChannels)
		logs = make([]Log, 0)
	}
}

func (r *Roga) consumeStdoutWrites(wg *sync.WaitGroup) {
	defer wg.Done()

	var flushed = false

	for {
		var operations []Operation

		var logs []Log

		if flushed {
			for {
				select {
				case writable := <-r.writingChannels.Stdout:
					collectWritable(writable, &operations, &logs)
				default:
					var stop = len(r.logQueue) == 0 && len(r.operationQueue) == 0

					if stop {
						close(r.writingChannels.Stdout)
						return
					}

					break
				}
			}
		}

		select {
		case flushed, _ = <-r.writingChannelsFlush.Stdout:
			if !flushed {
				flushed = true
				continue
			}
		default:
			if len(r.writingChannels.Stdout) < r.config.maxStdoutWriterChannelItems {
				continue
			}

			for i := 0; i < r.config.maxStdoutWriterChannelItems; i++ {
				collectWritable(<-r.writingChannels.Stdout, &operations, &logs)
			}
		}

		writeToStream("stdout", &operations, &logs, r)
	}
}

func (r *Roga) consumeFileWrites(wg *sync.WaitGroup) {
	defer wg.Done()

	var flushed = false

	for {
		var operations []Operation

		var logs []Log

		if flushed {
			for {
				select {
				case writable := <-r.writingChannels.File:
					collectWritable(writable, &operations, &logs)
				default:
					var stop = len(r.logQueue) == 0 && len(r.operationQueue) == 0

					if stop {
						close(r.writingChannels.File)
						return
					}

					break
				}
			}
		}

		select {
		case flushed, _ = <-r.writingChannelsFlush.File:
			if !flushed {
				flushed = true
				continue
			}
		default:
			if len(r.writingChannels.File) < r.config.maxFileWriterChannelItems {
				continue
			}
			for i := 0; i < r.config.maxFileWriterChannelItems; i++ {
				collectWritable(<-r.writingChannels.File, &operations, &logs)
			}
		}

		writeToStream("file", &operations, &logs, r)
	}
}

func (r *Roga) consumeExternalWrites(wg *sync.WaitGroup) {
	defer wg.Done()

	var flushed = false

	for {
		var operations []Operation

		var logs []Log

		if flushed {
			for {
				select {
				case writable := <-r.writingChannels.External:
					collectWritable(writable, &operations, &logs)
				default:
					var stop = len(r.logQueue) == 0 && len(r.operationQueue) == 0

					if stop {
						close(r.writingChannels.External)
						return
					}

					break
				}
			}
		}

		select {
		case flushed, _ = <-r.writingChannelsFlush.External:
			if !flushed {
				flushed = true
				continue
			}
		default:
			if len(r.writingChannels.External) < r.config.maxExternalWriterChannelItems {
				continue
			}

			for i := 0; i < r.config.maxExternalWriterChannelItems; i++ {
				collectWritable(<-r.writingChannels.External, &operations, &logs)
			}
		}

		writeToStream("external", &operations, &logs, r)
	}
}

func (r *Roga) consumeChannels() {
	go r.monitorAndUpdateSystemMetrics()

	go r.consumeProductionChannel()

	go r.consumeLogQueue()

	go r.consumeOperationQueue()

	go func() {
		r.consumptionSync.Add(1)
		defer r.consumptionSync.Done()

		var wg sync.WaitGroup

		for i := 0; i < r.config.maxStdoutWriters; i++ {
			go r.consumeStdoutWrites(&wg)
		}

		wg.Wait()
	}()

	go func() {
		r.consumptionSync.Add(1)
		defer r.consumptionSync.Done()

		var wg sync.WaitGroup

		for i := 0; i < r.config.maxFileWriters; i++ {
			go r.consumeFileWrites(&wg)
		}

		wg.Wait()
	}()

	go func() {
		r.consumptionSync.Add(1)
		defer r.consumptionSync.Done()

		var wg sync.WaitGroup

		for i := 0; i < r.config.maxExternalWriters; i++ {
			go r.consumeExternalWrites(&wg)
		}

		wg.Wait()
	}()
}

func addWritableToQueue(writable Writable, r *Roga) {
	if operation, ok := writable.(Operation); ok {
		var _, alreadyInBuffer = r.operationBuffer[operation.Id]

		if !alreadyInBuffer {
			return
		}

		r.operationBuffer[operation.Id] = operation
		r.operationQueue <- operation.Id
	} else if log, ok := writable.(Log); ok {
		r.logBuffer[log.Id] = log
		r.logQueue <- log.Id
	}
}

func collectWritable(writable Writable, operations *[]Operation, logs *[]Log) {
	if operation, ok := writable.(Operation); ok {
		*operations = append(*operations, operation)
	} else if log, ok := writable.(Log); ok {
		*logs = append(*logs, log)
	}
}

func writeToStream(stream string, operations *[]Operation, logs *[]Log, r *Roga) {
	var hasOperations = len(*operations) > 0
	var hasLogs = len(*logs) > 0

	switch stream {
	case "stdout":
		if hasOperations {
			r.writer.WriteOperationsToStdout(*operations, r)
		}

		if hasLogs {
			r.writer.WriteLogsToStdout(*logs, r)
		}
	case "file":
		if hasOperations {
			var file, cleanupFunc, err = getLogFileDescriptor(r, true)

			if err == nil {
				r.writer.WriteOperationsToFile(*operations, file, r)

				cleanupFunc(file)
			}
		}

		if hasLogs {
			var file, cleanupFunc, err = getLogFileDescriptor(r)

			if err == nil {
				r.writer.WriteLogsToFile(*logs, file, r)

				cleanupFunc(file)
			}
		}
	case "external":
		if hasOperations {
			r.writer.WriteOperationsToExternal(*operations, r)
		}

		if hasLogs {
			r.writer.WriteLogsToExternal(*logs, r)
		}
	}

	if hasOperations {
		*operations = make([]Operation, 0)
	}

	if hasLogs {
		*logs = make([]Log, 0)
	}
}
