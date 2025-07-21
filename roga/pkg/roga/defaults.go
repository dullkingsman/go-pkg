package roga

import (
	"time"
)

const (
	DefaultProductionQueueSize   int = 1024
	DefaultStdoutWriterQueueSize int = 1024
	DefaultFileWriterQueueSize   int = 1024

	DefaultSystemStatsCheckInterval = 5 * time.Second

	DefaultIdleChannelFlushInterval = 500 * time.Millisecond

	DefaultFileWriterBasePath            = "./logs/"
	DefaultFileLogsDirectoryGranularity  = time.Hour
	DefaultFileLogsDirectoryFormatLayout = "2006-01-02_15-04-05"

	DefaultWriteToStdout = true
	DefaultWriteToFile   = false

	DefaultOperationsFileName = "operations.jsonl"
	DefaultLogsFileName       = "logs.jsonl"
)
