package core

import (
	"github.com/google/uuid"
	"os"
	"sync"
	"time"
)

type (
	Roga struct {
		buffers         buffers
		channels        channels
		metricsLock     *sync.RWMutex
		consumptionSync *sync.WaitGroup
		started         bool
		lastWrite       time.Time
		rootOperation   Operation
		context         Context
		producer        Producer
		monitor         Monitor
		dispatcher      Dispatcher
		writer          Writer
		config          InstanceConfig
	}

	Config struct {
		ServiceName string
		Instance    *InstanceConfig
		Producer    Producer
		Monitor     Monitor
		Dispatcher  Dispatcher
		Writer      Writer
	}

	InstanceConfig struct {
		maxOperationQueueSize         int
		maxLogQueueSize               int
		maxProductionChannelItems     int
		maxStdoutWriterChannelItems   int
		maxFileWriterChannelItems     int
		maxExternalWriterChannelItems int
		maxStdoutWriters              int
		maxFileWriters                int
		maxExternalWriters            int
		systemStatsCheckInterval      time.Duration // in seconds
		writeToStdout                 bool
		writeToFile                   bool
		writeToExternal               bool
		fileWriterBasePath            string
		fileLogsDirectoryGranularity  time.Duration
		fileLogsDirectoryFormatLayout string // time format layout
		operationsFileName            string
		logsFileName                  string
	}

	Producer interface {
		LogFatal(args LogArgs, operation *Operation, context Context, framesToSkip int, ch *chan Writable) *Log
		LogError(args LogArgs, operation *Operation, context Context, framesToSkip int, ch *chan Writable) *Log
		LogInfo(args LogArgs, operation *Operation, context Context, framesToSkip int, ch *chan Writable) *Log
		LogWarn(args LogArgs, operation *Operation, context Context, framesToSkip int, ch *chan Writable) *Log
		LogDebug(args LogArgs, operation *Operation, context Context, framesToSkip int, ch *chan Writable) *Log
		BeginOperation(args OperationArgs, parent *Operation, measurementInitiator *MeasurementHandler, ch *chan Writable) *Operation
		EndOperation(operation *Operation, measurementFinalizer *MeasurementHandler, ch *chan Writable)
	}

	Monitor interface {
		GetCPUUsage() (usage float64, err error)
		GetMemoryStats() (total, free uint64, err error)
		GetSwapStats() (total, free uint64, err error)
		GetDiskStats(path string) (total, free uint64, err error)
	}

	Dispatcher interface {
		AddToOperationQueue(operations []Operation, queue *chan<- uuid.UUID)
		AddToLogQueue(logs []Log, queue *chan<- uuid.UUID)
		DispatchOperations(
			operations []Operation,
			channels *writingChannels,
		) []uuid.UUID // return the ids of ones that were not dispatched
		DispatchLogs(
			logs []Log,
			channels *writingChannels,
		) []uuid.UUID // return the ids of ones that were not dispatched
	}

	Writer interface {
		WriteOperationsToStdout(items []Operation, r *Roga)
		WriteOperationsToFile(items []Operation, file *os.File, r *Roga)
		WriteOperationsToExternal(items []Operation, r *Roga)
		WriteLogsToStdout(items []Log, r *Roga)
		WriteLogsToFile(items []Log, file *os.File, r *Roga)
		WriteLogsToExternal(items []Log, r *Roga)
	}

	Writable interface {
		String() string
	}

	MeasurementHandler func(*map[string]float64)

	buffers struct {
		operations map[uuid.UUID]Operation
		logs       map[uuid.UUID]Log
	}

	channels struct {
		operational channelGroup
		flush       actionChannelGroup
		stop        actionChannelGroup
	}

	channelGroup struct {
		production chan Writable
		queue      queueChannels
		writing    writingChannels
	}

	queueChannels struct {
		operation chan uuid.UUID
		log       chan uuid.UUID
	}

	writingChannels struct {
		stdout   chan Writable
		file     chan Writable
		external chan Writable
	}

	actionChannelGroup struct {
		production chan bool
		queue      queueActionChannels
		writing    writingActionChannels
	}

	queueActionChannels struct {
		operation chan bool
		log       chan bool
	}

	writingActionChannels struct {
		stdout   chan bool
		file     chan bool
		external chan bool
	}

	OperationArgs struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Actor       Actor   `json:"actor"`
	}

	LogArgs struct {
		Priority       *Priority       `json:"priority"`
		VerbosityClass *VerbosityClass `json:"verbosityClass"`
		Message        string          `json:"message"`
		Actor          Actor           `json:"actor"`
		Data           *interface{}    `json:"data"`
	}

	MonitorConfig struct {
		Interval int // in seconds
	}

	Replay struct {
		Id           uuid.UUID          `json:"id"`
		Name         string             `json:"name"`
		Index        int                `json:"index"` // the index of the replay in the list of replays for the same operation
		OperationId  *uuid.UUID         `json:"operationId"`
		Measurements map[string]float64 `json:"measurements"`
	}

	Operation struct {
		Writable
		r                     *Roga
		Id                    uuid.UUID             `json:"id"`
		Name                  string                `json:"name"`
		Description           *string               `json:"description"`
		BaseOperationId       *uuid.UUID            `json:"baseOperationId"`
		ParentId              *uuid.UUID            `json:"parentId"`
		OperationChildren     []uuid.UUID           `json:"operationChildren"`
		LogChildren           []uuid.UUID           `json:"logChildren"`
		EssentialMeasurements EssentialMeasurements `json:"essentialMeasurements"`
		Measurements          map[string]float64    `json:"measurements"`
		Actor                 Actor                 `json:"actor"`
	}

	Log struct {
		Writable
		Id             uuid.UUID      `json:"id"`
		Level          Level          `json:"level"`
		Priority       Priority       `json:"priority"`
		VerbosityClass VerbosityClass `json:"verbosityClass"`
		Message        string         `json:"message"`
		TracingId      uuid.UUID      `json:"tracingId"`
		OperationId    uuid.UUID      `json:"operationId"`
		Timestamp      time.Time      `json:"timestamp"`
		Stack          StackTrace     `json:"stack"`
		Actor          Actor          `json:"actor"`
		Context        Context        `json:"context"`
		Data           *interface{}   `json:"data"`
	}

	EssentialMeasurements struct {
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
	}

	Actor struct {
		Type           ActorType       `json:"type"`
		User           *User           `json:"user,omitempty"`
		ExternalSystem *ExternalSystem `json:"externalSystem,omitempty"`
	}

	User struct {
		Identifier      string  `json:"identifier"` // anything specific that can identify the user. E.g. if the user is not yet created the phone number and if they are, the id.
		Id              *string `json:"id,omitempty"`
		IdType          *string `json:"idType,omitempty"`
		SessionId       *string `json:"sessionId,omitempty"`
		SessionIdType   *string `json:"sessionIdType,omitempty"`
		Role            *string `json:"role,omitempty"`
		PermissionLevel *string `json:"permissionLevel,omitempty"`
		Type            *string `json:"type,omitempty"`
		PhoneNumber     *string `json:"phoneNumber,omitempty"`
		Email           *string `json:"email,omitempty"`
	}

	ExternalSystem struct {
		Id   string `json:"id"`
		Ip   string `json:"ip"`
		Name string `json:"name,omitempty"`
	}

	StackTrace struct {
		Crashed bool         `json:"crashed"`
		Frames  []StackFrame `json:"frames"`
	}

	StackFrame struct {
		File       string `json:"file"`
		Function   string `json:"function"`
		LineNumber int    `json:"lineNumber"`
	}

	Context struct {
		Environment          Environment          `json:"environment"`
		SystemSpecifications SystemSpecifications `json:"systemSpecifications"`
	}

	Environment struct {
		ApplicationEnvironment ApplicationEnvironment `json:"applicationEnvironment"`
		SystemEnvironment      SystemEnvironment      `json:"systemEnvironment"`
	}

	ApplicationEnvironment struct {
		ServiceName string `json:"serviceName"`
		Version     string `json:"version"`
		Env         string `json:"env"`
		Lang        string `json:"lang"`
		LangVersion string `json:"langVersion"`
		ProcessId   int    `json:"processId"`
	}

	SystemEnvironment struct {
		CpuUsage        float64 `json:"cpuUsage"`
		AvailableMemory uint64  `json:"availableMemory"`
		AvailableDisk   uint64  `json:"availableDisk"`
		AvailableSwap   uint64  `json:"availableSwap"`
	}

	SystemSpecifications struct {
		Os       string `json:"os"`
		Arch     string `json:"arch"`
		CpuCores int    `json:"cpuCores"`
		Memory   uint64 `json:"memory"`
		SwapSize uint64 `json:"swapSize"`
		DiskSize uint64 `json:"diskSize"`
		PageSize int    `json:"pageSize"`
	}

	ActorType      uint
	Priority       uint
	VerbosityClass uint
	Level          int
)
