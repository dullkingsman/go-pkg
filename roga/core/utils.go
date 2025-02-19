package core

import (
	"errors"
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
