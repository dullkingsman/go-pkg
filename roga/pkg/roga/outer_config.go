package roga

import (
	"github.com/dullkingsman/go-pkg/roga/internal/_map"
	"time"
)

type OuterInstanceConfig struct {
	VerbosityClass                *VerbosityClass
	LogLevel                      *Level
	ProductionBufferSize          *int
	WritingBufferSize             _map.NestedMap[WriteStream, EntryType, int]
	WriteBlock                    _map.NestedMap[WriteStream, EntryType, bool]
	StaleEntriesFlushInterval     *time.Duration
	SystemMetricsCheckInterval    *time.Duration // in seconds
	WriteToStdout                 *bool
	WriteToFile                   *bool
	WriteToExternal               *bool
	FileWriterBasePath            *string
	FileLogsDirectoryGranularity  *time.Duration
	FileLogsDirectoryFormatLayout *string // time format layout
	OperationsFileName            *string
	LogsFileName                  *string
	DontLogLoggerLogs             *bool
}

func (oic *OuterInstanceConfig) Inner() InnerInstanceConfig {
	var _config = DefaultInstanceConfig

	if oic == nil {
		return _config
	}

	if oic.VerbosityClass != nil {
		_config.verbosityClass = *oic.VerbosityClass
	}

	if oic.LogLevel != nil {
		_config.logLevel = *oic.LogLevel
	}
	if oic.ProductionBufferSize != nil {
		_config.productionBufferSize = *oic.ProductionBufferSize
	}

	// TODO write functions to patch these cleanly
	if oic.WritingBufferSize != nil {
		_config.writingBufferSize = oic.WritingBufferSize
	}
	if oic.WriteBlock != nil {
		_config.writeBlock = oic.WriteBlock
	}

	if oic.StaleEntriesFlushInterval != nil {
		_config.staleEntriesFlushInterval = *oic.StaleEntriesFlushInterval
	}
	if oic.SystemMetricsCheckInterval != nil {
		_config.systemMetricsCheckInterval = *oic.SystemMetricsCheckInterval
	}

	if oic.WriteToStdout != nil {
		_config.writeToStdout = *oic.WriteToStdout
	}

	if oic.WriteToFile != nil {
		_config.writeToFile = *oic.WriteToFile
	}

	if oic.FileWriterBasePath != nil {
		_config.fileWriterBasePath = *oic.FileWriterBasePath
	}

	if oic.FileLogsDirectoryGranularity != nil {
		_config.fileLogsDirectoryGranularity = *oic.FileLogsDirectoryGranularity
	}

	if oic.FileLogsDirectoryFormatLayout != nil {
		_config.fileLogsDirectoryFormatLayout = *oic.FileLogsDirectoryFormatLayout
	}

	if oic.OperationsFileName != nil {
		_config.operationsFileName = *oic.OperationsFileName
	}

	if oic.LogsFileName != nil {
		_config.logsFileName = *oic.LogsFileName
	}

	if oic.DontLogLoggerLogs != nil {
		_config.dontLogLoggerLogs = *oic.DontLogLoggerLogs
	}

	return _config
}
