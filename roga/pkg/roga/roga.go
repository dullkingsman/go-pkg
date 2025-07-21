package roga

import (
	"github.com/dullkingsman/go-pkg/roga/internal"
	"github.com/dullkingsman/go-pkg/roga/internal/file"
	"github.com/dullkingsman/go-pkg/roga/internal/queue"
	"github.com/dullkingsman/go-pkg/roga/internal/ticker"
	"github.com/dullkingsman/go-pkg/roga/internal/writer"
	"github.com/dullkingsman/go-pkg/roga/pkg/model"
	"github.com/dullkingsman/go-pkg/roga/writable"
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/google/uuid"
	"os"
	"sync"
	"time"
)

type Roga struct {
	systemMetricsMonitor *ticker.ManagedTicker
	staleEntriesMonitor  *ticker.ManagedTicker
	channels             Channels
	stdoutWriter         *writer.QueuedWriter[writable.Writable]
	fileWriter           *writer.QueuedWriter[writable.Writable]
	productionQueue      *queue.SelfConsumingQueue[writable.Writable]
	metricsLock          *sync.RWMutex
	wg                   *sync.WaitGroup
	started              bool
	lastWriteLock        *sync.RWMutex
	lastWrite            time.Time
	rootOperation        Operation
	context              model.Context
	currentSystemMetrics model.SystemMetrics
	stdoutFormatter      writable.Formatter
	producer             Producer
	monitor              Monitor
	writer               Writer
	config               InnerInstanceConfig
}

type (
	EntryType   int
	WriteStream int
)

const (
	EntryTypeOperation EntryType = iota
	EntryTypeLog
	EntryTypeEvent
	EntryTypeAudit
)

const (
	WriteStreamStdout WriteStream = iota
	WriteStreamFile
)

var (
	EntryTypeValues = []EntryType{
		EntryTypeOperation,
		EntryTypeLog,
		EntryTypeEvent,
		EntryTypeAudit,
	}

	WriteStreamValues = []WriteStream{
		WriteStreamStdout,
		WriteStreamFile,
	}

	EntryTypeName = map[EntryType]string{
		EntryTypeOperation: "operation",
		EntryTypeLog:       "log",
		EntryTypeEvent:     "event",
		EntryTypeAudit:     "audit",
	}

	WriteStreamName = map[WriteStream]string{
		WriteStreamStdout: "stdout",
		WriteStreamFile:   "file",
	}
)

func Init(cfg ...Config) Roga {
	if os.Geteuid() != 0 {
		utils.LogFatal("roga", "needs to run as root")
	}

	var _cfg = defaultConfig.FromIncoming(cfg...)

	var instanceConfig = _cfg.InstanceConfig.Inner()

	var r = Roga{
		context:         model.DefaultContext,
		config:          instanceConfig,
		metricsLock:     &sync.RWMutex{},
		lastWriteLock:   &sync.RWMutex{},
		wg:              &sync.WaitGroup{},
		started:         false,
		stdoutFormatter: _cfg.StdoutFormatter,
		producer:        _cfg.Producer,
		monitor:         _cfg.Monitor,
		writer:          _cfg.Writer,
		rootOperation: Operation{
			Id:          uuid.New(),
			Name:        "root",
			Description: utils.PtrOf("A program run!"),
			EssentialMeasurements: model.EssentialMeasurements{
				StartTime: time.Now().UTC(),
			},
			Actor: model.Actor{Type: 1},
		},
		channels: Channels{
			operational: ChannelGroup{
				production: make(chan writable.Writable, instanceConfig.productionBufferSize),
				writing: WriteStreamChannels{
					WriteStreamStdout: {
						EntryTypeOperation: make(chan writable.Writable, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeOperation]),
						EntryTypeLog:       make(chan writable.Writable, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeLog]),
						EntryTypeEvent:     make(chan writable.Writable, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeEvent]),
						EntryTypeAudit:     make(chan writable.Writable, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeAudit]),
					},
					WriteStreamFile: {
						EntryTypeOperation: make(chan writable.Writable, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeOperation]),
						EntryTypeLog:       make(chan writable.Writable, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeLog]),
						EntryTypeEvent:     make(chan writable.Writable, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeEvent]),
						EntryTypeAudit:     make(chan writable.Writable, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeAudit]),
					},
				},
			},
			flush: ActionChannelGroup{
				production: make(chan bool),
				writing: WriteStreamActionChannels{
					WriteStreamStdout: {
						EntryTypeOperation: make(chan bool, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeOperation]),
						EntryTypeLog:       make(chan bool, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeLog]),
						EntryTypeEvent:     make(chan bool, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeEvent]),
						EntryTypeAudit:     make(chan bool, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeAudit]),
					},
					WriteStreamFile: {
						EntryTypeOperation: make(chan bool, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeOperation]),
						EntryTypeLog:       make(chan bool, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeLog]),
						EntryTypeEvent:     make(chan bool, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeEvent]),
						EntryTypeAudit:     make(chan bool, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeAudit]),
					},
				},
			},
			stop: ActionChannelGroup{
				production: make(chan bool),
				writing: WriteStreamActionChannels{
					WriteStreamStdout: {
						EntryTypeOperation: make(chan bool, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeOperation]),
						EntryTypeLog:       make(chan bool, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeLog]),
						EntryTypeEvent:     make(chan bool, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeEvent]),
						EntryTypeAudit:     make(chan bool, instanceConfig.writingBufferSize[WriteStreamStdout][EntryTypeAudit]),
					},
					WriteStreamFile: {
						EntryTypeOperation: make(chan bool, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeOperation]),
						EntryTypeLog:       make(chan bool, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeLog]),
						EntryTypeEvent:     make(chan bool, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeEvent]),
						EntryTypeAudit:     make(chan bool, instanceConfig.writingBufferSize[WriteStreamFile][EntryTypeAudit]),
					},
				},
			},
		},
	}

	r.setApplicationSpecifics(_cfg)

	r.setCurrentSystemMetrics()

	r.rootOperation.RegisterInstance(&r)

	r.systemMetricsMonitor = ticker.NewManagedTicker(
		instanceConfig.GetSystemMetricsCheckInterval(),
		r.setCurrentSystemMetrics,
		nil,
	)

	r.staleEntriesMonitor = ticker.NewManagedTicker(
		instanceConfig.GetStaleEntriesFlushInterval(),
		r.flushStaleEntries,
		nil,
	)

	r.productionQueue = queue.NewQueue(
		r.dispatchWritableObjects,
		r.config.productionBufferSize,
		nil,
		nil,
		nil,
	)

	r.stdoutWriter = writer.NewQueuedWriter[writable.Writable](
		func(currentWriter *writer.SyncedWriter) *writer.SyncedWriter {
			if currentWriter == nil {
				return writer.NewSyncedWriter(os.Stdout, 256*1024)
			}

			return currentWriter
		},
		DefaultStdoutFormatter{},
		DefaultStdoutWriterQueueSize,
		[]*queue.SelfConsumingQueue[writable.Writable]{r.productionQueue},
		nil,
	)

	r.fileWriter = writer.NewQueuedWriter[writable.Writable](
		func(currentWriter *writer.SyncedWriter) *writer.SyncedWriter {
			if currentWriter == nil {
				var _file, err = file.GetLogFileDescriptor(
					r.config.GetFileWriterBasePath(),
					r.config.GetCurrentLogDirName(),
					r.config.GetLogFileName(EntryTypeLog),
				)

				if err != nil {
					println(err.Error())
				}

				return writer.NewSyncedWriter(_file, 256*1024)
			}

			return currentWriter
		},
		DefaultStdoutFormatter{},
		DefaultStdoutWriterQueueSize,
		[]*queue.SelfConsumingQueue[writable.Writable]{r.productionQueue},
		nil,
	)

	r.productionQueue.AddDependents(r.stdoutWriter.GetQueue(), r.fileWriter.GetQueue())

	r.started = true

	r.dispatchWritableObject(r.rootOperation)

	return r
}

func (r *Roga) Start() {
	return
	if !r.started {
		r.started = true

		r.dispatchWritableObject(r.rootOperation)

		r.startConsuming()

		utils.LogInfo("roga:startup", utils.CyanString("started"))
	}
}

func (r *Roga) Flush() {
	//if r.started {
	//r.channels.flush.production <- true
	r.productionQueue.Flush()
	//}
}

func (r *Roga) Stop() {
	if r.started {
		r.started = false

		r.rootOperation.EndOperation()

		r.fileWriter.GetQueue().Stop()
		r.stdoutWriter.GetQueue().Stop()
		r.productionQueue.Stop()
		//r.stopConsuming()

		r.wg.Wait()

		//r.stopMonitoring()

		utils.LogInfo("roga:cleanup", utils.CyanString("stopped"))
	}
}

func (r *Roga) LogFatal(args LogArgs) {
	produceLog(
		r.producer,
		LevelFatal,
		TypeNormal,
		args,
		&r.rootOperation,
		&r.currentSystemMetrics,
		1,
		//r.channels.operational.production,
		r.productionQueue,
	)

	os.Exit(1)
}

func (r *Roga) LogError(args LogArgs) {
	produceLog(
		r.producer,
		LevelError,
		TypeNormal,
		args,
		&r.rootOperation,
		&r.currentSystemMetrics,
		1,
		//r.channels.operational.production,
		r.productionQueue,
	)
}

func (r *Roga) LogWarn(args LogArgs) {
	produceLog(
		r.producer,
		LevelWarn,
		TypeNormal,
		args,
		&r.rootOperation,
		&r.currentSystemMetrics,
		1,
		//r.channels.operational.production,
		r.productionQueue,
	)
}

func (r *Roga) LogInfo(args LogArgs) {
	produceLog(
		r.producer,
		LevelInfo,
		TypeNormal,
		args,
		&r.rootOperation,
		&r.currentSystemMetrics,
		1,
		//r.channels.operational.production,
		r.productionQueue,
	)
}

func (r *Roga) LogDebug(args LogArgs) {
	produceLog(
		r.producer,
		LevelDebug,
		TypeNormal,
		args,
		&r.rootOperation,
		&r.currentSystemMetrics,
		1,
		//r.channels.operational.production,
		r.productionQueue,
	)
}

func (r *Roga) AuditAction(args AuditLogArgs) {
	produceLog(
		r.producer,
		LevelInfo,
		TypeAudit,
		args.LogArgs,
		&r.rootOperation,
		&r.currentSystemMetrics,
		1,
		//r.channels.operational.production,
		r.productionQueue,
	)
}

func (r *Roga) CaptureEvent(args EventLogArgs) {
	produceLog(
		r.producer,
		LevelInfo,
		TypeEvent,
		args.LogArgs,
		&r.rootOperation,
		&r.currentSystemMetrics,
		1,
		//r.channels.operational.production,
		r.productionQueue,
	)
}

func (r *Roga) BeginOperation(args OperationArgs, measurementInitiator ...MeasurementHandler) *Operation {
	return beginOperation(
		r.producer,
		args,
		&r.rootOperation,
		&r.context,
		//r.channels.operational.production,
		r.productionQueue,
		measurementInitiator...,
	)
}

func (r *Roga) stopMonitoring() {
	r.staleEntriesMonitor.Stop()
	r.systemMetricsMonitor.Stop()
	//r.systemMetricsMonitorControls.stop <- true
	//close(r.systemMetricsMonitorControls.stop)
	//
	//r.staleEntriesMonitorControls.stop <- true
	//close(r.staleEntriesMonitorControls.stop)
}

func (r *Roga) startConsuming() {
	r.consumeProduction()

	for _, entryType := range EntryTypeValues {
		for _, writeStream := range WriteStreamValues {
			if writeStream == WriteStreamStdout && r.config.writeToStdout {
				r.consumeWrites(writeStream, entryType)
			}

			if writeStream == WriteStreamFile && r.config.writeToFile {
				r.consumeWrites(writeStream, entryType)
			}

			//if r.config.writeToExternal {
			//	for i := 0; i < r.config.noOfExternalWriters[entryType]; i++ {
			//		r.consumeWrites(writeStream, entryType)
			//	}
			//}
		}
	}
}

func (r *Roga) stopConsuming() {
	r.channels.stop.production <- true
	close(r.channels.operational.production)
	close(r.channels.flush.production) // TODO maybe I shouldn't close these to prevent race conditions causing panic when the monitor tires to flush
	close(r.channels.stop.production)
}

func (r *Roga) monitorAndFlushStaleEntries() {

	//r.wg.Add(1)
	//internal.monitor(
	//	"stale entries",
	//	r.staleEntriesMonitorControls.stop,
	//	r.staleEntriesMonitorControls.pause,
	//	r.staleEntriesMonitorControls.resume,
	//	r.config.staleEntriesFlushInterval,
	//	r.wg,
	//	func() {
	//		if r.lastWrite.Before(time.Now().Add(-r.config.staleEntriesFlushInterval)) {
	//			r.Flush()
	//		}
	//	},
	//	r.config.dontLogLoggerLogs,
	//)
}

func (r *Roga) monitorAndUpdateSystemMetrics() {
	r.systemMetricsMonitor.Stop()
	//r.wg.Add(1)
	//internal.monitor(
	//	"system metrics",
	//	r.systemMetricsMonitorControls.stop,
	//	r.systemMetricsMonitorControls.pause,
	//	r.systemMetricsMonitorControls.resume,
	//	r.config.systemMetricsCheckInterval,
	//	r.wg,
	//	r.setCurrentSystemMetrics,
	//	r.config.dontLogLoggerLogs,
	//)
}

func (r *Roga) consumeProduction() {
	r.wg.Add(1)
	go internal.ConsumeQueue(
		"production",
		r.channels.operational.production,
		r.channels.stop.production,
		r.channels.flush.production,
		[]<-chan writable.Writable{},
		[]chan<- bool{
			r.channels.flush.writing[WriteStreamStdout][EntryTypeOperation],
			r.channels.flush.writing[WriteStreamStdout][EntryTypeLog],
			r.channels.flush.writing[WriteStreamStdout][EntryTypeEvent],
			r.channels.flush.writing[WriteStreamStdout][EntryTypeAudit],
			r.channels.flush.writing[WriteStreamFile][EntryTypeOperation],
			r.channels.flush.writing[WriteStreamFile][EntryTypeLog],
			r.channels.flush.writing[WriteStreamFile][EntryTypeEvent],
			r.channels.flush.writing[WriteStreamFile][EntryTypeAudit],
		},
		[]chan<- bool{
			r.channels.stop.writing[WriteStreamStdout][EntryTypeOperation],
			r.channels.stop.writing[WriteStreamStdout][EntryTypeLog],
			r.channels.stop.writing[WriteStreamStdout][EntryTypeEvent],
			r.channels.stop.writing[WriteStreamStdout][EntryTypeAudit],
			r.channels.stop.writing[WriteStreamFile][EntryTypeOperation],
			r.channels.stop.writing[WriteStreamFile][EntryTypeLog],
			r.channels.stop.writing[WriteStreamFile][EntryTypeEvent],
			r.channels.stop.writing[WriteStreamFile][EntryTypeAudit],
		},
		r.wg,
		r.dispatchWritableObjects,
		r.config.dontLogLoggerLogs,
	)
}

func (r *Roga) consumeWrites(writeStream WriteStream, entryType EntryType) {
	if writeStream == WriteStreamStdout {
		return
	}

	r.wg.Add(1)
	go internal.ConsumeQueue(
		WriteStreamName[writeStream]+" "+EntryTypeName[entryType]+" writes",
		r.channels.operational.writing[writeStream][entryType],
		r.channels.stop.writing[writeStream][entryType],
		r.channels.flush.writing[writeStream][entryType],
		[]<-chan writable.Writable{
			r.channels.operational.production,
		},
		nil,
		nil,
		r.wg,
		func(items []writable.Writable) {
			r.writeToStream(items, writeStream, entryType, r.stdoutFormatter)

			r.lastWriteLock.Lock()
			defer r.lastWriteLock.Unlock()

			r.lastWrite = time.Now().UTC()
		},
		r.config.dontLogLoggerLogs,
	)
}

func (r *Roga) setApplicationSpecifics(cfg Config) {
	r.context.Application.SetApplicationSpecifics(
		cfg.Code,
		cfg.Version,
		cfg.Env,
		cfg.Node,
		cfg.Instance,
		cfg.Name,
	)
}

func (r *Roga) setCurrentSystemMetrics() {
	r.metricsLock.Lock()
	defer r.metricsLock.Unlock()

	var (
		cpuUsage, cpuErr                   = r.monitor.GetCPUUsage()
		totalMemory, freeMemory, memoryErr = r.monitor.GetMemoryStats()
		totalSwap, freeSwap, swapErr       = r.monitor.GetSwapStats()
		totalDisk, freeDisk, diskErr       = r.monitor.GetDiskStats("/")
	)

	if cpuErr == nil {
		r.currentSystemMetrics.CpuUsage = cpuUsage
	}

	if memoryErr == nil {
		r.context.System.Memory = totalMemory
		r.currentSystemMetrics.AvailableMemory = freeMemory
	}

	if swapErr == nil {
		r.context.System.SwapSize = totalSwap
		r.currentSystemMetrics.AvailableSwap = freeSwap
	}

	if diskErr == nil {
		r.context.System.DiskSize = totalDisk
		r.currentSystemMetrics.AvailableDisk = freeDisk
	}
}

func (r *Roga) flushStaleEntries() {
	if r.lastWrite.Before(time.Now().Add(-r.config.GetStaleEntriesFlushInterval())) {
		//r.Flush()
	}
}

func (r *Roga) dispatchWritableObjects(items []writable.Writable) {
	r.stdoutWriter.GetQueue().EnqueueMany(items)
	r.fileWriter.GetQueue().EnqueueMany(items)

	//for _, writable := range items {
	//	r.dispatchWritableObject(writable)
	//}
}

func (r *Roga) dispatchWritableObject(writable writable.Writable) {
	if r.config.writeToStdout {
		r.dispatchWritable(writable, WriteStreamStdout)
	}

	if r.config.writeToFile {
		r.dispatchWritable(writable, WriteStreamFile)
	}
}

func (r *Roga) dispatchWritable(writable writable.Writable, stream WriteStream) {
	var channel = r.channels.operational.writing[stream][EntryTypeOperation]

	var log = utils.SafeCastValue[Log](writable)

	if log != nil {
		if log.Type == TypeAudit {
			channel = r.channels.operational.writing[stream][EntryTypeAudit]
		} else if log.Type == TypeEvent {
			channel = r.channels.operational.writing[stream][EntryTypeEvent]
		} else {
			channel = r.channels.operational.writing[stream][EntryTypeLog]
		}
	}

	// TODO CONSIDER WHAT HAPPENS IF THE CHANNEL IS CLOSED
	channel <- writable
}

func (r *Roga) writeToStream(
	items []writable.Writable,
	stream WriteStream,
	entryType EntryType,
	stdoutFormatter writable.Formatter,
) {
	switch stream {
	case WriteStreamStdout:
		var syncedWriter = writer.NewSyncedWriter(os.Stdout, 256*1024)

		var _, err = r.writer.WriteToStdout(items, entryType, syncedWriter, stdoutFormatter)

		_ = syncedWriter.Flush()

		if err != nil {
			// TODO what happens here
		}

	case WriteStreamFile:

		var _file, err = file.GetLogFileDescriptor(
			r.config.GetFileWriterBasePath(),
			r.config.GetCurrentLogDirName(),
			r.config.GetLogFileName(entryType),
		)

		if err != nil {
			// TODO what happens here
			break
		}

		var syncedWriter = writer.NewSyncedWriter(_file, 512*1024)

		_, err = r.writer.WriteToFile(items, entryType, syncedWriter)

		_ = syncedWriter.Flush()

		if err != nil {
			// TODO what happens here
		}

		_ = _file.Close()
	}
}
