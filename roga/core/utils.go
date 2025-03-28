package core

import (
	"errors"
	"github.com/dullkingsman/go-pkg/utils"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// getCPUUsage retrieves system-wide CPU usage (Linux/macOS only)
func getCPUUsage() (float64, error) {
	var idle1, total1, err = readCPUStats()

	if err != nil {
		return 0, err
	}

	time.Sleep(500 * time.Millisecond)

	idle2, total2, err := readCPUStats()

	if err != nil {
		return 0, err
	}

	var (
		idleTicks  = float64(idle2 - idle1)
		totalTicks = float64(total2 - total1)
	)

	return (1.0 - (idleTicks / totalTicks)) * 100.0, nil
}

// getMemoryStats retrieves total and free memory (Linux/macOS only)
func getMemoryStats() (total, free uint64, err error) {
	var info syscall.Sysinfo_t

	err = syscall.Sysinfo(&info)

	if err != nil {
		return 0, 0, err
	}

	total = info.Totalram * uint64(info.Unit)

	free = info.Freeram * uint64(info.Unit)

	return
}

// getDiskStats retrieves total and free disk space (Linux/macOS only)
func getDiskStats(path string) (total, free uint64, err error) {
	var stat syscall.Statfs_t

	err = syscall.Statfs(path, &stat)

	if err != nil {
		return
	}

	total = stat.Blocks * uint64(stat.Bsize)

	free = stat.Bfree * uint64(stat.Bsize)

	return
}

// getSwapStats retrieves total and free swap memory (Linux only)
func getSwapStats() (total, free uint64, err error) {
	return readSwapStats()
}

// readCPUStats reads CPU statistics from /proc/stat (Linux/macOS)
func readCPUStats() (idle, total uint64, err error) {
	file, err := os.ReadFile("/proc/stat")

	if err != nil {
		return
	}

	var lines = strings.Split(string(file), "\n")

	if len(lines) < 1 {
		err = errors.New("could not get stats")

		return
	}

	for _, line := range lines {
		var fields = strings.Fields(line)

		if len(fields) < 8 || fields[0] != "cpu" {
			continue
		}

		var values []uint64

		for _, field := range fields[1:] {
			var val, _err = strconv.ParseUint(field, 10, 64)

			err = _err

			if err != nil {
				return
			}

			values = append(values, val)
		}

		// Assign total and idle times
		idle = values[3] // Idle time is the 4th field

		for _, v := range values {
			total += v
		}

		return
	}

	err = errors.New("could not get stats")

	return
}

// readSwapStats retrieves total and free swap memory (Linux only)
func readSwapStats() (total, free uint64, err error) {
	data, err := os.ReadFile("/proc/meminfo")

	if err != nil {
		return
	}

	var lines = strings.Split(string(data), "\n")

	if len(lines) < 2 {
		err = errors.New("could not get swap stats")

		return
	}

	for _, line := range lines {
		var fields = strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		var key, value = fields[0], fields[1]

		v, err := strconv.ParseUint(value, 10, 64)

		if err != nil {
			continue
		}

		if key == "SwapTotal:" {
			total = v
		} else if key == "SwapFree:" {
			free = v
		}
	}

	return
}

func getStackFrames(framesToSkip int) []StackFrame {
	var (
		stack = make([]StackFrame, 0)
		pc    []uintptr
	)

	var retrieved = runtime.Callers(framesToSkip+1, pc)

	if retrieved == 0 {
		return nil
	}

	var frames = runtime.CallersFrames(pc)

	for {
		var frame, ok = frames.Next()

		if !ok {
			break
		}

		if frame.File == "" && frame.Line == 0 && frame.Function == "" {
			continue
		}

		stack = append(stack, StackFrame{
			File:       frame.File,
			Function:   frame.Function,
			LineNumber: frame.Line,
		})
	}

	return stack
}

func (a LogArgs) ToLog() Log {
	var log = Log{
		Message: a.Message,
		Actor:   a.Actor,
		Data:    a.Data,
	}

	if a.VerbosityClass != nil {
		log.VerbosityClass = *a.VerbosityClass
	}

	if a.Priority != nil {
		log.Priority = *a.Priority
	}

	return log
}

func (a OperationArgs) ToOperation() Operation {
	return Operation{
		Name:        a.Name,
		Description: a.Description,
		Actor:       a.Actor,
	}
}

// getCurrentTimeRoundedTo rounds the current time to the nearest given interval (in nanoseconds)
func getCurrentTimeRoundedTo(interval time.Duration) time.Time {
	now := time.Now()

	// Calculate the nearest multiple of the interval
	unixTime := now.UnixNano()                                        // Convert time to nanoseconds
	roundedTimeNano := (unixTime / int64(interval)) * int64(interval) // Round down

	// If past halfway, round up
	if unixTime%int64(interval) >= int64(interval)/2 {
		roundedTimeNano += int64(interval)
	}

	// Convert nanoseconds back to time.Time
	roundedTime := time.Unix(0, roundedTimeNano)

	return roundedTime
}

// getLogFileDescriptor returns a file descriptor for the given file name in the log base directory
func getLogFileDescriptor(r *Roga, operations ...bool) (file *os.File, cleanupFunc func(file *os.File), err error) {
	var isOperations = false

	if len(operations) > 0 {
		isOperations = operations[0]
	}

	var logsBasePath = r.config.fileWriterBasePath +
		getCurrentTimeRoundedTo(
			r.config.fileLogsDirectoryGranularity,
		).UTC().Format(
			r.config.fileLogsDirectoryFormatLayout,
		)

	err = os.MkdirAll(logsBasePath, os.ModePerm)

	if err != nil {
		utils.LogError("roga:file-descriptor", err.Error())
	}

	var filePath = logsBasePath + "/"

	if isOperations {
		filePath += r.config.operationsFileName
	} else {
		filePath += r.config.logsFileName
	}

	file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil || file == nil {
		if err == nil {
			err = errors.New("could not open file " + utils.GreyString(filePath))
		}

		utils.LogError("roga:file-descriptor", err.Error())

		return
	}

	cleanupFunc = func(file *os.File) {
		var err = file.Close()

		if err != nil {
			utils.LogError("roga:file-descriptor", err.Error())
		}
	}

	return
}

func addWritableToQueue(writable Writable, r *Roga) {
	if operation, ok := writable.(Operation); ok {
		var _, alreadyInBuffer = r.buffers.operations[operation.Id]

		if !alreadyInBuffer {
			return
		}

		r.buffers.operations[operation.Id] = operation
		r.channels.operational.queue.operation <- operation.Id
	} else if log, ok := writable.(Log); ok {
		r.buffers.logs[log.Id] = log
		r.channels.operational.queue.log <- log.Id
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

	r.lastWriteLock.Lock()

	r.lastWrite = time.Now().UTC()

	r.lastWriteLock.Unlock()
}
