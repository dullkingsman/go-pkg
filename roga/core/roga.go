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
		monitorControls: monitorControls{
			stop:   make(chan bool),
			pause:  make(chan bool),
			resume: make(chan bool),
		},
		buffers: buffers{
			operations: make(map[uuid.UUID]Operation),
			logs:       make(map[uuid.UUID]Log),
		},
		channels: channels{
			operational: channelGroup{
				production: make(chan Writable, _config.Instance.maxProductionChannelItems),
				queue: queueChannels{
					operation: make(chan uuid.UUID, _config.Instance.maxOperationQueueSize),
					log:       make(chan uuid.UUID, _config.Instance.maxLogQueueSize),
				},
				writing: writingChannels{
					stdout:   make(chan Writable, _config.Instance.maxStdoutWriterChannelItems),
					file:     make(chan Writable, _config.Instance.maxFileWriterChannelItems),
					external: make(chan Writable, _config.Instance.maxExternalWriterChannelItems),
				},
			},
			flush: actionChannelGroup{
				production: make(chan bool),
				queue: queueActionChannels{
					operation: make(chan bool),
					log:       make(chan bool),
				},
				writing: writingActionChannels{
					stdout:   make(chan bool),
					file:     make(chan bool),
					external: make(chan bool),
				},
			},
			stop: actionChannelGroup{
				production: make(chan bool),
				queue: queueActionChannels{
					operation: make(chan bool),
					log:       make(chan bool),
				},
				writing: writingActionChannels{
					stdout:   make(chan bool),
					file:     make(chan bool),
					external: make(chan bool),
				},
			},
		},
	}

	instance.context.Environment.ApplicationEnvironment.ServiceName = _config.ServiceName

	instance.rootOperation.r = &instance

	return instance
}

func (r *Roga) StopMonitoring() {
	r.ResumeMonitoring()
	r.monitorControls.stop <- true
}

func (r *Roga) PauseMonitoring() {
	r.monitorControls.pause <- true
}

func (r *Roga) ResumeMonitoring() {
	r.monitorControls.resume <- true
}

func (r *Roga) Start() {
	if !r.started {
		r.started = true
		r.consumeChannels()
	}
}

func (r *Roga) Flush() {
	r.channels.flush.writing.stdout <- true
	r.channels.flush.writing.file <- true
	r.channels.flush.writing.external <- true

	r.channels.flush.queue.operation <- true
	r.channels.flush.queue.log <- true

	r.channels.flush.production <- true
}

func (r *Roga) Stop() {
	r.StopMonitoring()

	r.channels.stop.writing.stdout <- true
	r.channels.stop.writing.file <- true
	r.channels.stop.writing.external <- true

	r.channels.stop.queue.operation <- true
	r.channels.stop.queue.log <- true

	r.channels.stop.production <- true
}

func (r *Roga) Recover() {
	r.Stop()

	r.Wait()

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
		&r.channels.operational.production,
	)
}

func (r *Roga) LogError(args LogArgs) {
	r.producer.LogError(
		args,
		&r.rootOperation,
		r.context,
		1,
		&r.channels.operational.production,
	)
}

func (r *Roga) LogWarn(args LogArgs) {
	r.producer.LogWarn(
		args,
		&r.rootOperation,
		r.context,
		1,
		&r.channels.operational.production,
	)
}

func (r *Roga) LogInfo(args LogArgs) {
	r.producer.LogInfo(
		args,
		&r.rootOperation,
		r.context,
		1,
		&r.channels.operational.production,
	)
}

func (r *Roga) LogDebug(args LogArgs) {
	r.producer.LogDebug(
		args,
		&r.rootOperation,
		r.context,
		1,
		&r.channels.operational.production,
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
		&r.channels.operational.production,
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
		&o.r.channels.operational.production,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogError(args LogArgs) {
	var log = o.r.producer.LogError(
		args,
		o,
		o.r.context,
		1,
		&o.r.channels.operational.production,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogWarn(args LogArgs) {
	var log = o.r.producer.LogWarn(
		args,
		o,
		o.r.context,
		1,
		&o.r.channels.operational.production,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogInfo(args LogArgs) {
	var log = o.r.producer.LogInfo(
		args,
		o,
		o.r.context,
		1,
		&o.r.channels.operational.production,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogDebug(args LogArgs) {
	var log = o.r.producer.LogDebug(
		args,
		o,
		o.r.context,
		1,
		&o.r.channels.operational.production,
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
		&o.r.channels.operational.production,
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
		&o.r.channels.operational.production,
	)
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

func (r *Roga) consumeProductionChannel() {
	r.consumptionSync.Add(1)
	defer r.consumptionSync.Done()

	var stopped = false

	for {
		if stopped {
			for {
				select {
				case writable := <-r.channels.operational.production:
					addWritableToQueue(writable, r)
				default:
					close(r.channels.operational.production)

					r.channels.stop.queue.operation <- true
					r.channels.stop.queue.log <- true

					return
				}
			}
		}

		select {
		case <-r.channels.stop.production:
			if !stopped {
				stopped = true
				continue
			}
		case <-r.channels.flush.production:
			for {
				select {
				case writable := <-r.channels.operational.production:
					addWritableToQueue(writable, r)
				default:
					r.channels.flush.queue.operation <- true
					r.channels.flush.queue.log <- true

					break
				}
			}
		default:
			if len(r.channels.operational.production) < r.config.maxProductionChannelItems {
				continue
			}

			for i := 0; i < r.config.maxProductionChannelItems; i++ {
				addWritableToQueue(<-r.channels.operational.production, r)
			}
		}
	}
}

func (r *Roga) consumeOperationQueue() {
	r.consumptionSync.Add(1)
	defer r.consumptionSync.Done()

	var stopped = false

	for {
		var operations []Operation

		if stopped {
			for {
				select {
				case operationId := <-r.channels.operational.queue.operation:
					var operation, ok = r.buffers.operations[operationId]

					if !ok {
						continue
					}

					operations = append(operations, operation)
				default:
					var stop = len(r.channels.operational.production) == 0

					r.channels.stop.writing.stdout <- true
					r.channels.stop.writing.file <- true
					r.channels.stop.writing.external <- true

					if stop {
						close(r.channels.operational.queue.operation)
						return
					}

					break
				}
			}
		}

		select {
		case <-r.channels.stop.queue.operation:
			if !stopped {
				stopped = true
				continue
			}
		case <-r.channels.flush.queue.operation:
			for {
				select {
				case operationId := <-r.channels.operational.queue.operation:
					var operation, ok = r.buffers.operations[operationId]

					if !ok {
						continue
					}

					operations = append(operations, operation)
				default:
					r.channels.flush.writing.stdout <- true
					r.channels.flush.writing.file <- true
					r.channels.flush.writing.external <- true

					break
				}
			}
		default:
			if len(r.channels.operational.queue.operation) < r.config.maxOperationQueueSize {
				continue
			}

			operations = make([]Operation, r.config.maxOperationQueueSize)

			for i := 0; i < r.config.maxOperationQueueSize; i++ {
				var operationId = <-r.channels.operational.queue.operation

				var operation, ok = r.buffers.operations[operationId]

				if !ok {
					continue
				}

				operations[i] = operation
			}
		}

		r.dispatcher.DispatchOperations(operations, &r.channels.operational.writing)
		operations = make([]Operation, 0)
	}
}

func (r *Roga) consumeLogQueue() {
	r.consumptionSync.Add(1)
	defer r.consumptionSync.Done()

	var stopped = false

	for {
		var logs []Log

		if stopped {
			for {
				select {
				case logId := <-r.channels.operational.queue.log:
					var log, ok = r.buffers.logs[logId]

					if !ok {
						continue
					}

					logs = append(logs, log)
				default:
					var stop = len(r.channels.operational.production) == 0

					r.channels.stop.writing.stdout <- true
					r.channels.stop.writing.file <- true
					r.channels.stop.writing.external <- true

					if stop {
						close(r.channels.operational.queue.log)
						return
					}

					break
				}
			}
		}

		select {
		case <-r.channels.stop.queue.log:
			if !stopped {
				stopped = true
				continue
			}
		case <-r.channels.flush.queue.log:
			for {
				select {
				case logId := <-r.channels.operational.queue.log:
					var log, ok = r.buffers.logs[logId]

					if !ok {
						continue
					}

					logs = append(logs, log)
				default:
					r.channels.flush.writing.stdout <- true
					r.channels.flush.writing.file <- true
					r.channels.flush.writing.external <- true

					break
				}
			}
		default:
			if len(r.channels.operational.queue.log) < r.config.maxLogQueueSize {
				continue
			}

			logs = make([]Log, r.config.maxLogQueueSize)

			for i := 0; i < r.config.maxLogQueueSize; i++ {
				var logId = <-r.channels.operational.queue.log

				var log, ok = r.buffers.logs[logId]

				if !ok {
					continue
				}

				logs[i] = log
			}
		}

		r.dispatcher.DispatchLogs(logs, &r.channels.operational.writing)
		logs = make([]Log, 0)
	}
}

func (r *Roga) consumeStdoutWrites(wg *sync.WaitGroup) {
	defer wg.Done()

	var stopped = false

	for {
		var operations []Operation

		var logs []Log

		if stopped {
			for {
				select {
				case writable := <-r.channels.operational.writing.stdout:
					collectWritable(writable, &operations, &logs)
				default:
					var stop = len(r.channels.operational.queue.log) == 0 && len(r.channels.operational.queue.operation) == 0

					if stop {
						close(r.channels.operational.writing.stdout)
						return
					}

					break
				}
			}
		}

		select {
		case <-r.channels.stop.writing.stdout:
			if !stopped {
				stopped = true
				continue
			}
		case <-r.channels.flush.writing.stdout:
			for {
				select {
				case writable := <-r.channels.operational.writing.stdout:
					collectWritable(writable, &operations, &logs)
				default:
					break
				}
			}
		default:
			if len(r.channels.operational.writing.stdout) < r.config.maxStdoutWriterChannelItems {
				continue
			}

			for i := 0; i < r.config.maxStdoutWriterChannelItems; i++ {
				collectWritable(<-r.channels.operational.writing.stdout, &operations, &logs)
			}
		}

		writeToStream("stdout", &operations, &logs, r)
	}
}

func (r *Roga) consumeFileWrites(wg *sync.WaitGroup) {
	defer wg.Done()

	var stopped = false

	for {
		var operations []Operation

		var logs []Log

		if stopped {
			for {
				select {
				case writable := <-r.channels.operational.writing.file:
					collectWritable(writable, &operations, &logs)
				default:
					var stop = len(r.channels.operational.queue.log) == 0 && len(r.channels.operational.queue.operation) == 0

					if stop {
						close(r.channels.operational.writing.file)
						return
					}

					break
				}
			}
		}

		select {
		case <-r.channels.stop.writing.file:
			if !stopped {
				stopped = true
				continue
			}
		case <-r.channels.flush.writing.file:
			for {
				select {
				case writable := <-r.channels.operational.writing.file:
					collectWritable(writable, &operations, &logs)
				default:
					break
				}
			}
		default:
			if len(r.channels.operational.writing.file) < r.config.maxFileWriterChannelItems {
				continue
			}
			for i := 0; i < r.config.maxFileWriterChannelItems; i++ {
				collectWritable(<-r.channels.operational.writing.file, &operations, &logs)
			}
		}

		writeToStream("file", &operations, &logs, r)
	}
}

func (r *Roga) consumeExternalWrites(wg *sync.WaitGroup) {
	defer wg.Done()

	var stopped = false

	for {
		var operations []Operation

		var logs []Log

		if stopped {
			for {
				select {
				case writable := <-r.channels.operational.writing.external:
					collectWritable(writable, &operations, &logs)
				default:
					var stop = len(r.channels.operational.queue.log) == 0 && len(r.channels.operational.queue.operation) == 0

					if stop {
						close(r.channels.operational.writing.external)
						return
					}

					break
				}
			}
		}

		select {
		case <-r.channels.stop.writing.external:
			if !stopped {
				stopped = true
				continue
			}
		case <-r.channels.flush.writing.external:
			for {
				select {
				case writable := <-r.channels.operational.writing.external:
					collectWritable(writable, &operations, &logs)
				default:
					break
				}
			}
		default:
			if len(r.channels.operational.writing.external) < r.config.maxExternalWriterChannelItems {
				continue
			}

			for i := 0; i < r.config.maxExternalWriterChannelItems; i++ {
				collectWritable(<-r.channels.operational.writing.external, &operations, &logs)
			}
		}

		writeToStream("external", &operations, &logs, r)
	}
}
