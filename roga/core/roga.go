package roga

import (
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/google/uuid"
	"os"
	"sync"
	"time"
)

func Init(config ...Config) Roga {
	if os.Geteuid() != 0 {
		utils.LogFatal("roga", "needs to run as root")
	}

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
		lastWriteLock:   &sync.RWMutex{},
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
		metricMonitorControls: monitorControls{
			stop:   make(chan bool, 1),
			pause:  make(chan bool, 1),
			resume: make(chan bool, 1),
		},
		idleChannelMonitorControls: monitorControls{
			stop:   make(chan bool, 1),
			pause:  make(chan bool, 1),
			resume: make(chan bool, 1),
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

	instance.context.Application.Name = _config.Name
	instance.context.Application.Code = _config.Code

	var (
		product    = getProductIdentifier()
		provider   = getCloudProvider(product)
		instanceId = getCloudInstanceID(provider)
		machineId  = getMachineID()
		macAddress = getMacAddress()
	)

	instance.context.System.Product = product
	instance.context.System.InstanceId = instanceId
	instance.context.System.MachineId = machineId
	instance.context.System.MacAddress = macAddress

	instance.rootOperation.r = &instance

	return instance
}

func (r *Roga) StopSystemMonitoring() {
	r.ResumeSystemMonitoring()
	r.metricMonitorControls.stop <- true

	utils.LogInfo("roga:cleanup", "stopped system monitoring")
}

func (r *Roga) PauseSystemMonitoring() {
	r.metricMonitorControls.pause <- true

	utils.LogInfo("roga:cleanup", "paused system monitoring")
}

func (r *Roga) ResumeSystemMonitoring() {
	r.metricMonitorControls.resume <- true

	utils.LogInfo("roga:cleanup", "resumed system monitoring")
}

func (r *Roga) StopIdleChannelMonitoring() {
	r.ResumeIdleChannelMonitoring()

	r.idleChannelMonitorControls.stop <- true

	utils.LogInfo("roga:cleanup", "stopped idle channel monitoring")
}

func (r *Roga) PauseIdleChannelMonitoring() {
	r.idleChannelMonitorControls.pause <- true

	utils.LogInfo("roga:cleanup", "paused idle channel monitoring")
}

func (r *Roga) ResumeIdleChannelMonitoring() {
	r.idleChannelMonitorControls.resume <- true

	utils.LogInfo("roga:cleanup", "resumed idle channel monitoring")
}

func (r *Roga) Start() {
	if !r.started {
		r.started = true

		r.startingMonitoring()

		r.startConsuming()
	}

	utils.LogInfo("roga:startup", "started")
}

func (r *Roga) Flush() {
	r.channels.flush.writing.stdout <- true
	r.channels.flush.writing.file <- true
	r.channels.flush.writing.external <- true

	utils.LogInfo("roga:cleanup:consumption", "flushed writes")

	r.channels.flush.queue.operation <- true
	r.channels.flush.queue.log <- true

	utils.LogInfo("roga:cleanup:consumption", "flushed queues")

	r.channels.flush.production <- true

	utils.LogInfo("roga:cleanup:consumption", "flushed production")
}

func (r *Roga) StopConsuming() {
	r.channels.stop.production <- true

	utils.LogInfo("roga:cleanup:consumption", "stopped consuming production")

	r.channels.stop.queue.operation <- true
	r.channels.stop.queue.log <- true

	utils.LogInfo("roga:cleanup:consumption", "stopped consuming queues")

	r.channels.stop.writing.stdout <- true
	r.channels.stop.writing.file <- true
	r.channels.stop.writing.external <- true

	utils.LogInfo("roga:cleanup:consumption", "stopped consuming writes")
}

func (r *Roga) Stop(flush ...bool) {
	if len(flush) > 0 && flush[0] {
		r.Flush()
	}

	r.StopIdleChannelMonitoring()

	r.StopSystemMonitoring()

	r.StopConsuming()

	utils.LogInfo("roga:cleanup", "stopped")
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
		r.currentSystemMetrics,
		1,
		&r.channels.operational.production,
	)
}

func (r *Roga) LogError(args LogArgs) {
	r.producer.LogError(
		args,
		&r.rootOperation,
		r.currentSystemMetrics,
		1,
		&r.channels.operational.production,
	)
}

func (r *Roga) LogWarn(args LogArgs) {
	r.producer.LogWarn(
		args,
		&r.rootOperation,
		r.currentSystemMetrics,
		1,
		&r.channels.operational.production,
	)
}

func (r *Roga) LogInfo(args LogArgs) {
	r.producer.LogInfo(
		args,
		&r.rootOperation,
		r.currentSystemMetrics,
		1,
		&r.channels.operational.production,
	)
}

func (r *Roga) LogDebug(args LogArgs) {
	r.producer.LogDebug(
		args,
		&r.rootOperation,
		r.currentSystemMetrics,
		1,
		&r.channels.operational.production,
	)
}

func (r *Roga) AuditAction(args AuditLogArgs) {
	r.producer.AuditAction(
		args,
		&r.rootOperation,
		1,
		&r.channels.operational.production,
	)
}

func (r *Roga) CaptureEvent(args EventLogArgs) {
	r.producer.CaptureEvent(
		args,
		&r.rootOperation,
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
		&r.context,
		_measurementInitiator,
		&r.channels.operational.production,
	)

	operation.r = r

	return operation
}

func (o *Operation) LogFatal(args LogArgs) {
	var log = o.r.producer.LogFatal(
		args,
		o,
		o.r.currentSystemMetrics,
		1,
		&o.r.channels.operational.production,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogError(args LogArgs) {
	var log = o.r.producer.LogError(
		args,
		o,
		o.r.currentSystemMetrics,
		1,
		&o.r.channels.operational.production,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogWarn(args LogArgs) {
	var log = o.r.producer.LogWarn(
		args,
		o,
		o.r.currentSystemMetrics,
		1,
		&o.r.channels.operational.production,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogInfo(args LogArgs) {
	var log = o.r.producer.LogInfo(
		args,
		o,
		o.r.currentSystemMetrics,
		1,
		&o.r.channels.operational.production,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) LogDebug(args LogArgs) {
	var log = o.r.producer.LogDebug(
		args,
		o,
		o.r.currentSystemMetrics,
		1,
		&o.r.channels.operational.production,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) AuditAction(args AuditLogArgs) {
	var log = o.r.producer.AuditAction(
		args,
		o,
		1,
		&o.r.channels.operational.production,
	)

	o.LogChildren = append(o.LogChildren, log.Id)
}

func (o *Operation) CaptureEvent(args EventLogArgs) {
	var log = o.r.producer.CaptureEvent(
		args,
		o,
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
		&o.r.context,
		_measurementInitiator,
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
		_measurementFinalizer,
		&o.r.channels.operational.production,
	)
}

func (r *Roga) startingMonitoring() {
	SetCurrentSystemMetrics(r)

	go r.monitorAndFlushIdleChannels()

	go r.monitorAndUpdateSystemMetrics()
}

func (r *Roga) startConsuming() {
	go r.consumeProductionChannel()

	go r.consumeLogQueue()

	go r.consumeOperationQueue()

	go func() {
		r.consumptionSync.Add(1)
		defer r.consumptionSync.Done()

		var wg sync.WaitGroup

		for i := 0; i < 1; /*r.config.maxStdoutWriters \/*TODO*\/*/ i++ {
			go r.consumeStdoutWrites(&wg, i)
		}

		wg.Wait()
	}()

	go func() {
		r.consumptionSync.Add(1)
		defer r.consumptionSync.Done()

		var wg sync.WaitGroup

		for i := 0; i < 1; /*r.config.maxFileWriters \/*TODO*\/*/ i++ {
			go r.consumeFileWrites(&wg, i)
		}

		wg.Wait()
	}()

	go func() {
		r.consumptionSync.Add(1)
		defer r.consumptionSync.Done()

		var wg sync.WaitGroup

		for i := 0; i < 1; /*r.config.maxExternalWriters \/*TODO*\/*/ i++ {
			go r.consumeExternalWrites(&wg, i)
		}

		wg.Wait()
	}()
}

func (r *Roga) consumeProductionChannel() {
	r.consumptionSync.Add(1)
	defer r.consumptionSync.Done()

	utils.LogInfo("roga:startup", "consuming production...")

	var stopped = false

	for {
		time.Sleep(100 * time.Millisecond)
		if stopped {
			for {
				select {
				case writable := <-r.channels.operational.production:
					addWritableToQueue(writable, r)
				default:
					close(r.channels.operational.production)

					//r.channels.stop.queue.operation <- true
					//r.channels.stop.queue.log <- true

					utils.LogInfo("roga:cleanup", "stopped production consumption")

					return
				}
			}
		}

		select {
		case <-r.channels.stop.production:
			if !stopped {
				stopped = true
				utils.LogInfo("roga:ops", "signaled production consumption to stop ")

				continue
			}
		case <-r.channels.flush.production:
		FlushProductionBreak:
			for {
				select {
				case writable := <-r.channels.operational.production:
					addWritableToQueue(writable, r)
				default:
					r.channels.flush.queue.operation <- true
					r.channels.flush.queue.log <- true
					break FlushProductionBreak
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

	utils.LogInfo("roga:startup", "consuming operations...")

	var stopped = false

	for {
		time.Sleep(100 * time.Millisecond)
		var operations []Operation

		if stopped {
		StopBreakStatement:
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

					//r.channels.stop.writing.stdout <- true
					//r.channels.stop.writing.file <- true
					//r.channels.stop.writing.external <- true

					if stop {
						close(r.channels.operational.queue.operation)

						utils.LogInfo("roga:cleanup", "stopped operation consumption")

						return
					}

					utils.LogInfo("roga:cleanup", "skipped stopping operation consumption")

					break StopBreakStatement
				}
			}
		}

		select {
		case <-r.channels.stop.queue.operation:
			if !stopped {
				stopped = true

				utils.LogInfo("roga:ops", "signaled operation consumption to stop")

				continue
			}
		case <-r.channels.flush.queue.operation:
		QueueOperationFlushBreak:
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

					break QueueOperationFlushBreak
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

	utils.LogInfo("roga:startup", "consuming logs...")

	var stopped = false

	for {
		time.Sleep(100 * time.Millisecond)
		var logs []Log

		if stopped {
		QueueOperationStopBreak:
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

					//r.channels.stop.writing.stdout <- true
					//r.channels.stop.writing.file <- true
					//r.channels.stop.writing.external <- true

					if stop {
						close(r.channels.operational.queue.log)

						utils.LogInfo("roga:cleanup", "stopped log consumption")

						return
					}

					utils.LogInfo("roga:cleanup", "skipped stopping log consumption")

					break QueueOperationStopBreak
				}
			}
		}

		select {
		case <-r.channels.stop.queue.log:
			if !stopped {
				stopped = true

				utils.LogInfo("roga:ops", "signaled log consumption to stop ")

				continue
			}
		case <-r.channels.flush.queue.log:
		FlushQueueBreak:
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

					break FlushQueueBreak
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

func (r *Roga) consumeStdoutWrites(wg *sync.WaitGroup, index int) {
	wg.Add(1)
	defer wg.Done()

	utils.LogInfo("roga:startup", "consuming stdout (%d)...", index)

	var stopped = false

	for {
		time.Sleep(100 * time.Millisecond)

		var operations []Operation

		var logs []Log

		if stopped {
		CheckDrainedBeforeFinalStop:
			for {
				select {
				case writable := <-r.channels.operational.writing.stdout:
					collectWritable(writable, &operations, &logs)
				default:
					var stop = len(r.channels.operational.queue.log) == 0 && len(r.channels.operational.queue.operation) == 0

					if stop {
						close(r.channels.operational.writing.stdout)

						utils.LogInfo("roga:cleanup", "stopped stdout")

						return
					}

					utils.LogInfo("roga:cleanup", "skipped stopping stdout")

					break CheckDrainedBeforeFinalStop
				}
			}
		}

		select {
		case <-r.channels.stop.writing.stdout:
			if !stopped {
				stopped = true

				utils.LogInfo("roga:ops", "signaled stdout writes to stop ")

				continue
			}
		case <-r.channels.flush.writing.stdout:
		FlushAllPendingItems:
			for {
				select {
				case writable := <-r.channels.operational.writing.stdout:
					collectWritable(writable, &operations, &logs)
				default:
					break FlushAllPendingItems
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

func (r *Roga) consumeFileWrites(wg *sync.WaitGroup, index int) {
	wg.Add(1)
	defer wg.Done()

	utils.LogInfo("roga:startup", "consuming file (%d)...", index)

	var stopped = false

	for {
		time.Sleep(100 * time.Millisecond)

		var operations []Operation

		var logs []Log

		if stopped {
		CheckDrainedBeforeFinalStop:
			for {
				select {
				case writable := <-r.channels.operational.writing.file:
					collectWritable(writable, &operations, &logs)
				default:
					var stop = len(r.channels.operational.queue.log) == 0 && len(r.channels.operational.queue.operation) == 0

					if stop {
						close(r.channels.operational.writing.file)

						utils.LogInfo("roga:cleanup", "stopped files")

						return
					}

					utils.LogInfo("roga:cleanup", "skipped stopping external")

					break CheckDrainedBeforeFinalStop
				}
			}
		}

		select {
		case <-r.channels.stop.writing.file:
			if !stopped {
				stopped = true

				utils.LogInfo("roga:ops", "signaled file writes to stop ")

				continue
			}
		case <-r.channels.flush.writing.file:
		FlushAllPendingItems:
			for {
				select {
				case writable := <-r.channels.operational.writing.file:
					collectWritable(writable, &operations, &logs)
				default:
					break FlushAllPendingItems
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

func (r *Roga) consumeExternalWrites(wg *sync.WaitGroup, index int) {
	wg.Add(1)
	defer wg.Done()

	utils.LogInfo("roga:startup", "consuming external (%d)...", index)

	var stopped = false

	for {
		time.Sleep(100 * time.Millisecond)

		var operations []Operation

		var logs []Log

		if stopped {
			for {
			CheckDrainedBeforeFinalStop:
				select {
				case writable := <-r.channels.operational.writing.external:
					collectWritable(writable, &operations, &logs)
				default:
					var stop = len(r.channels.operational.queue.log) == 0 && len(r.channels.operational.queue.operation) == 0

					if stop {
						close(r.channels.operational.writing.external)

						utils.LogInfo("roga:cleanup", "stopped external")

						return
					}

					utils.LogInfo("roga:cleanup", "skipped stopping external")

					break CheckDrainedBeforeFinalStop
				}
			}
		}

		select {
		case <-r.channels.stop.writing.external:
			if !stopped {
				stopped = true

				utils.LogInfo("roga:ops", "signaled external writes to stop")

				continue
			}
		case <-r.channels.flush.writing.external:
		FlushAllPendingItems:
			for {
				select {
				case writable := <-r.channels.operational.writing.external:
					collectWritable(writable, &operations, &logs)
				default:
					break FlushAllPendingItems
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
