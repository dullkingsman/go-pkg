package roga

import (
	"github.com/dullkingsman/go-pkg/roga/internal/_map"
	time2 "github.com/dullkingsman/go-pkg/roga/internal/utils"
	"time"
)

type InnerInstanceConfig struct {
	verbosityClass                VerbosityClass
	logLevel                      Level
	productionBufferSize          int
	writingBufferSize             _map.NestedMap[WriteStream, EntryType, int]
	writeBlock                    _map.NestedMap[WriteStream, EntryType, bool]
	staleEntriesFlushInterval     time.Duration
	systemMetricsCheckInterval    time.Duration // in seconds
	writeToStdout                 bool
	writeToFile                   bool
	fileWriterBasePath            string
	fileLogsDirectoryGranularity  time.Duration
	fileLogsDirectoryFormatLayout string // time format layout
	operationsFileName            string
	logsFileName                  string
	dontLogLoggerLogs             bool
}

func (ic InnerInstanceConfig) Outer() OuterInstanceConfig {
	return OuterInstanceConfig{
		VerbosityClass:                &ic.verbosityClass,
		LogLevel:                      &ic.logLevel,
		ProductionBufferSize:          &ic.productionBufferSize,
		WritingBufferSize:             ic.writingBufferSize,
		WriteBlock:                    ic.writeBlock,
		StaleEntriesFlushInterval:     &ic.staleEntriesFlushInterval,
		SystemMetricsCheckInterval:    &ic.systemMetricsCheckInterval,
		WriteToStdout:                 &ic.writeToStdout,
		WriteToFile:                   &ic.writeToFile,
		FileWriterBasePath:            &ic.fileWriterBasePath,
		FileLogsDirectoryGranularity:  &ic.fileLogsDirectoryGranularity,
		FileLogsDirectoryFormatLayout: &ic.fileLogsDirectoryFormatLayout,
		OperationsFileName:            &ic.operationsFileName,
		LogsFileName:                  &ic.logsFileName,
		DontLogLoggerLogs:             &ic.dontLogLoggerLogs,
	}
}

var DefaultInstanceConfig = InnerInstanceConfig{
	productionBufferSize: DefaultProductionQueueSize,
	writingBufferSize: _map.NestedMap[WriteStream, EntryType, int]{
		WriteStreamStdout: {
			EntryTypeOperation: DefaultStdoutWriterQueueSize,
			EntryTypeLog:       DefaultStdoutWriterQueueSize,
			EntryTypeEvent:     DefaultStdoutWriterQueueSize,
			EntryTypeAudit:     DefaultStdoutWriterQueueSize,
		},
		WriteStreamFile: {
			EntryTypeOperation: DefaultFileWriterQueueSize,
			EntryTypeLog:       DefaultFileWriterQueueSize,
			EntryTypeEvent:     DefaultFileWriterQueueSize,
			EntryTypeAudit:     DefaultFileWriterQueueSize,
		},
	},
	staleEntriesFlushInterval:     DefaultIdleChannelFlushInterval,
	systemMetricsCheckInterval:    DefaultSystemStatsCheckInterval,
	writeToStdout:                 DefaultWriteToStdout,
	writeToFile:                   DefaultWriteToFile,
	fileWriterBasePath:            DefaultFileWriterBasePath,
	fileLogsDirectoryGranularity:  DefaultFileLogsDirectoryGranularity,
	fileLogsDirectoryFormatLayout: DefaultFileLogsDirectoryFormatLayout,
	operationsFileName:            DefaultOperationsFileName,
	logsFileName:                  DefaultLogsFileName,
	dontLogLoggerLogs:             true,
}

func (ic InnerInstanceConfig) GetLogFileName(entryType EntryType) string {
	var filePath = ""

	switch entryType {
	case EntryTypeOperation:
		filePath += ic.operationsFileName
	case EntryTypeAudit:
		filePath += "audit." + ic.logsFileName
	case EntryTypeEvent:
		filePath += "event." + ic.logsFileName
	case EntryTypeLog:
		filePath += ic.logsFileName
	}

	return filePath
}

func (ic InnerInstanceConfig) GetCurrentLogDirName() string {
	return time2.GetTimeRoundedTo(
		ic.fileLogsDirectoryGranularity,
	).UTC().Format(
		ic.fileLogsDirectoryFormatLayout,
	)
}

func (ic InnerInstanceConfig) GetFileWriterBasePath() string {
	return ic.fileWriterBasePath
}

func (ic InnerInstanceConfig) GetStaleEntriesFlushInterval() time.Duration {
	return ic.staleEntriesFlushInterval
}

func (ic InnerInstanceConfig) GetSystemMetricsCheckInterval() time.Duration {
	return ic.systemMetricsCheckInterval
}
