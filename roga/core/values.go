package roga

import (
	"github.com/dullkingsman/go-pkg/utils"
	"os"
	"runtime"
	"time"
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

	ActorTypeSystem         ActorType = 0 // system
	ActorTypeUser           ActorType = 1 // user
	ActorTypeExternalSystem ActorType = 2 // external system

	DefaultLogQueueSize               int = 1000
	DefaultOperationQueueSize         int = 1000
	DefaultProductionChannelItems     int = 1000
	DefaultStdoutWriterChannelItems   int = 1000
	DefaultFileWriterChannelItems     int = 1000
	DefaultExternalWriterChannelItems int = 1000

	DefaultSystemStatsCheckInterval = 10 // in seconds

	DefaultIdleChannelFlushInterval = 10 // in seconds

	DefaultFileWriterBasePath            = "./logs"
	DefaultFileLogsDirectoryGranularity  = time.Hour
	DefaultFileLogsDirectoryFormatLayout = "2006-01-02_15-04-05"

	DefaultWriteToStdout   = true
	DefaultWriteToFile     = true
	DefaultWriteToExternal = true

	DefaultOperationsFileName               = "operations"
	DefaultLogsFileName                     = "logs"
	DefaultLogsFormat                       = ".bson"
	CloudProviderUnknown      CloudProvider = 0
	CloudProviderAWS          CloudProvider = 1
	CloudProviderGCP          CloudProvider = 2
	CloudProviderAzure        CloudProvider = 3
	CloudProviderVmware       CloudProvider = 5
	CloudProviderVirtualBox   CloudProvider = 6
	CloudProviderKvmQemu      CloudProvider = 7

	FileTypeBinary Type = 0
	FileTypeText   Type = 1
)

var (
	DefaultMaxStdoutWriters   = 2 * runtime.NumCPU()
	DefaultMaxFileWriters     = 2 * runtime.NumCPU()
	DefaultMaxExternalWriters = 4 * runtime.NumCPU()

	defaultOperationContext = Context{
		Application: Application{
			Lang:        "go",
			LangVersion: runtime.Version(),
			ProcessId:   os.Getpid(),
		},
		System: SystemSpecifications{
			Os:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			CpuCores: runtime.NumCPU(),
			PageSize: os.Getpagesize(),
		},
	}

	defaultRogaConfig = Config{
		Instance:                 utils.PtrOf(defaultInstanceConfig.Outer()),
		StdoutOperationFormatter: &DefaultOperationFormatter{},
		StdoutLogFormatter:       &DefaultLogFormatter{},
		Producer:                 &DefaultProducer{},
		Monitor:                  &DefaultMonitor{},
		Dispatcher:               &DefaultDispatcher{},
		Writer:                   &DefaultWriter{},
		FileFormat:               utils.PtrOf("binary"),
	}

	defaultInstanceConfig = InstanceConfig{
		maxOperationQueueSize:         DefaultOperationQueueSize,
		maxLogQueueSize:               DefaultLogQueueSize,
		maxProductionChannelItems:     DefaultProductionChannelItems,
		maxStdoutWriterChannelItems:   DefaultStdoutWriterChannelItems,
		maxFileWriterChannelItems:     DefaultFileWriterChannelItems,
		maxExternalWriterChannelItems: DefaultExternalWriterChannelItems,
		maxStdoutWriters:              DefaultMaxStdoutWriters,
		maxFileWriters:                DefaultMaxFileWriters,
		maxExternalWriters:            DefaultMaxExternalWriters,
		idleChannelFlushInterval:      DefaultIdleChannelFlushInterval,
		systemStatsCheckInterval:      DefaultSystemStatsCheckInterval,
		writeToStdout:                 DefaultWriteToStdout,
		writeToFile:                   DefaultWriteToFile,
		writeToExternal:               DefaultWriteToExternal,
		fileWriterBasePath:            DefaultFileWriterBasePath,
		fileLogsDirectoryGranularity:  DefaultFileLogsDirectoryGranularity,
		fileLogsDirectoryFormatLayout: DefaultFileLogsDirectoryFormatLayout,
		operationsFileName:            DefaultOperationsFileName,
		logsFileName:                  DefaultLogsFileName,
		logsFormat:                    DefaultLogsFormat,
	}
)
