package roga

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dullkingsman/go-pkg/utils"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// getCloudInstanceID retrieves the instance ID based on the detected cloud provider.
func getCloudInstanceID(provider CloudProvider) *string {
	var client = &http.Client{Timeout: 5 * time.Second} // Longer timeout for actual data retrieval

	var instanceID *string

	switch provider {
	case CloudProviderAWS:
		// For AWS, first get a token for IMDSv2
		var tokenReq, _ = http.NewRequest("PUT", "http://169.254.169.254/latest/api/token", nil)

		tokenReq.Header.Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")

		tokenResp, tokenErr := client.Do(tokenReq)

		if tokenErr != nil || tokenResp.StatusCode != http.StatusOK {
			fmt.Printf("Error getting AWS IMDSv2 token: %v\n", tokenErr)
			return nil
		}

		defer func(Body io.ReadCloser) {
			var err = Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}(tokenResp.Body)

		token, _ := io.ReadAll(tokenResp.Body)

		var req, _ = http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/instance-id", nil)

		req.Header.Set("X-aws-ec2-metadata-token", string(token))

		var resp, err = client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			defer func(Body io.ReadCloser) {
				var err = Body.Close()

				if err != nil {
					fmt.Println(err)
				}
			}(resp.Body)

			var body, _ = io.ReadAll(resp.Body)

			instanceID = utils.PtrOf(strings.TrimSpace(string(body)))
		} else {
			fmt.Printf("Error getting AWS instance ID: %v\n", err)
		}

	case CloudProviderGCP:
		var req, _ = http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/instance/id", nil)

		req.Header.Set("Metadata-Flavor", "Google")

		var resp, err = client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			defer func(Body io.ReadCloser) {
				var err = Body.Close()

				if err != nil {
					fmt.Println(err)
				}
			}(resp.Body)

			var body, _ = io.ReadAll(resp.Body)

			instanceID = utils.PtrOf(strings.TrimSpace(string(body)))
		} else {
			fmt.Printf("Error getting GCP instance ID: %v\n", err)
		}

	case CloudProviderAzure:
		var req, _ = http.NewRequest("GET", "http://169.254.169.254/metadata/instance/compute/vmId?api-version=2021-02-01", nil)
		req.Header.Set("Metadata", "true")
		var resp, err = client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			defer func(Body io.ReadCloser) {
				var err = Body.Close()

				if err != nil {
					fmt.Println(err)
				}
			}(resp.Body)

			var body, _ = io.ReadAll(resp.Body)

			instanceID = utils.PtrOf(strings.TrimSpace(string(body)))

		} else {
			fmt.Printf("Error getting Azure instance ID: %v\n", err)
		}
	}

	return instanceID
}

// getCloudProvider attempts to detect the current cloud provider.
// It checks for well-known metadata service IP addresses and headers.
func getCloudProvider(product ProductIdentifier) CloudProvider {
	if product.Name != nil && *product.Name != "" {
		var lowerProductName = strings.ToLower(*product.Name)

		switch true {
		case strings.Contains(lowerProductName, strings.ToLower("HVM domU")),
			strings.Contains(lowerProductName, strings.ToLower("Amazon")),
			strings.Contains(lowerProductName, strings.ToLower("EC2")),
			strings.Contains(lowerProductName, strings.ToLower("Amazon EC2")):
			return CloudProviderAWS
		case strings.Contains(lowerProductName, strings.ToLower("Google Compute Engine")),
			strings.Contains(lowerProductName, strings.ToLower("Google")):
			return CloudProviderGCP
		case strings.Contains(lowerProductName, strings.ToLower("Virtual Machine")),
			strings.Contains(lowerProductName, strings.ToLower("Azure")):
			return CloudProviderAzure
		case strings.Contains(lowerProductName, strings.ToLower("VMware Virtual Platform")),
			strings.Contains(lowerProductName, strings.ToLower("VMware")):
			return CloudProviderVmware
		case strings.Contains(lowerProductName, strings.ToLower("VirtualBox")):
			return CloudProviderVirtualBox
		case strings.Contains(lowerProductName, strings.ToLower("Standard PC")),
			strings.Contains(lowerProductName, strings.ToLower("KVM")),
			strings.Contains(lowerProductName, strings.ToLower("QEMU")),
			strings.Contains(lowerProductName, strings.ToLower("QEMU Virtual Machine")):
			return CloudProviderKvmQemu
		}
	}

	var client = &http.Client{Timeout: 2 * time.Second}

	var req, _ = http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/", nil)

	resp, err := client.Do(req)

	if err == nil {
		defer func(Body io.ReadCloser) {
			var err = Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)

		if resp.StatusCode == http.StatusOK && resp.Header.Get("Server") == "EC2" {
			return CloudProviderAWS
		}
	}

	req, _ = http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/", nil)

	req.Header.Set("Metadata-Flavor", "Google")

	resp, err = client.Do(req)

	if err == nil {
		defer func(Body io.ReadCloser) {
			var err = Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)

		if resp.StatusCode == http.StatusOK {
			return CloudProviderGCP
		}
	}

	req, _ = http.NewRequest("GET", "http://169.254.169.254/metadata/instance?api-version=2021-02-01", nil)

	req.Header.Set("Metadata", "true")

	resp, err = client.Do(req)

	if err == nil {
		defer func(Body io.ReadCloser) {
			var err = Body.Close()

			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)

		if resp.StatusCode == http.StatusOK {
			var body, _ = io.ReadAll(resp.Body)

			if bytes.Contains(body, []byte(`"azEnvironment"`)) { // A common field in Azure IMDS response
				return CloudProviderAzure
			}
		}
	}

	return CloudProviderUnknown
}

func getProductIdentifier() ProductIdentifier {
	var product = ProductIdentifier{}

	switch runtime.GOOS {
	case "linux":
		var name, err = os.ReadFile("/sys/class/dmi/id/product_name")

		if err != nil {
			fmt.Printf("error reading product name: %v\n", err)
		} else {
			product.Name = utils.PtrOf(string(name))
		}

		uuid, err := os.ReadFile("/sys/class/dmi/id/product_uuid")

		if err != nil {
			fmt.Printf("error reading product uuid: %v\n", err)
		} else {
			product.Uuid = utils.PtrOf(string(uuid))
		}

		serial, err := os.ReadFile("/sys/class/dmi/id/product_serial")

		if err != nil {
			fmt.Printf("error reading product serial: %v\n", err)
		} else {
			product.Serial = utils.PtrOf(string(serial))
		}
	default:
		fmt.Printf("Hardware UUID retrieval not implemented for %s\n", runtime.GOOS)
	}

	return product
}

// getMachineID attempts to retrieve the OS-level machine ID.
func getMachineID() *string {
	var machineID *string

	switch runtime.GOOS {
	case "linux":
		var content, err = os.ReadFile("/etc/machine-id")

		if err == nil {
			if len(content) > 0 {
				machineID = utils.PtrOf(strings.TrimSpace(string(content)))
			}
		} else {
			fmt.Printf("error reading /etc/machine-id: %v\n", err)
		}
	default:
		fmt.Printf("Machine ID retrieval not implemented for %s\n", runtime.GOOS)
	}

	return machineID
}

// getMacAddress retrieves the MAC address of the first non-loopback network interface.
func getMacAddress() *string {
	var interfaces, err = net.Interfaces()

	if err != nil {
		fmt.Printf("Error getting network interfaces: %v\n", err)
		return nil
	}

	for _, _interface := range interfaces {
		if _interface.Flags&net.FlagUp != 0 && _interface.Flags&net.FlagLoopback == 0 {
			var mac = _interface.HardwareAddr.String()

			if mac != "" {
				return &mac
			}
		}
	}

	return nil
}

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
		Actor:   utils.ValueOr(a.Actor, Actor{}),
		Data:    a.Data,
		Event:   a.Event,
		Outcome: a.Outcome,
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
		Actor:       utils.ValueOr(a.Actor, Actor{}),
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
func getLogFileDescriptor(r *Roga, operations ...bool) (normal *os.File, audit *os.File, event *os.File, operation *os.File, cleanupFunc func(file *os.File), err error) {
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

	var (
		normalFilePath    = logsBasePath + "/"
		auditFilePath     = normalFilePath
		eventFilePath     = normalFilePath
		operationFilePath = normalFilePath
	)

	if isOperations {
		operationFilePath += r.config.operationsFileName

		operation, err = os.OpenFile(operationFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil || operation == nil {
			if err == nil {
				err = errors.New("could not open operation file " + utils.GreyString(operationFilePath))
			}

			utils.LogError("roga:operation-file-descriptor", err.Error())

			return
		}
	} else {
		normalFilePath += "normal." + r.config.logsFileName
		auditFilePath += "audit." + r.config.logsFileName
		eventFilePath += "event." + r.config.logsFileName

		normal, err = os.OpenFile(normalFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil || normal == nil {
			if err == nil {
				err = errors.New("could not open normal file " + utils.GreyString(normalFilePath))
			}

			utils.LogError("roga:normal-file-descriptor", err.Error())

			return
		}

		audit, err = os.OpenFile(auditFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil || audit == nil {
			if err == nil {
				err = errors.New("could not open audit file " + utils.GreyString(auditFilePath))
			}

			utils.LogError("roga:audit-file-descriptor", err.Error())

			return
		}

		event, err = os.OpenFile(eventFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil || event == nil {
			if err == nil {
				err = errors.New("could not open event file " + utils.GreyString(eventFilePath))
			}

			utils.LogError("roga:event-file-descriptor", err.Error())

			return
		}
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
		var idSuffix = "-start"

		if !operation.EssentialMeasurements.EndTime.IsZero() {
			idSuffix = "-end"
		}

		r.buffers.operations.Write(operation.Id.String()+idSuffix, operation)

		r.channels.operational.queue.operation <- operation.Id.String() + idSuffix

		return
	}

	if log, ok := writable.(Log); ok {
		var item = r.buffers.logs.Read(log.Id)

		if item == nil {
			r.buffers.logs.Write(log.Id, log)
			r.channels.operational.queue.log <- log.Id
		}
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
	r.writeSync.Add(1)
	defer r.writeSync.Done()

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
			var _, _, _, operation, cleanupFunc, err = getLogFileDescriptor(r, true)

			if err == nil {
				r.writer.WriteOperationsToFile(*operations, operation, r)

				cleanupFunc(operation)
			}
		}

		if hasLogs {
			var normal, audit, event, _, cleanupFunc, err = getLogFileDescriptor(r)

			if err == nil {
				r.writer.WriteLogsToFile(*logs, normal, audit, event, r)

				cleanupFunc(normal)
				cleanupFunc(audit)
				cleanupFunc(event)
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
		for _, operation := range *operations {
			if !operation.EssentialMeasurements.EndTime.IsZero() {
				r.buffers.operations.Remove(operation.Id.String() + "-start")
				r.buffers.operations.Remove(operation.Id.String() + "-end")
			}
		}

		*operations = make([]Operation, 0)
	}

	if hasLogs {
		for _, log := range *logs {
			r.buffers.logs.Remove(log.Id)
		}

		*logs = make([]Log, 0)
	}

	r.lastWriteLock.Lock()

	r.lastWrite = time.Now().UTC()

	r.lastWriteLock.Unlock()
}

func (b *buffer[T, H]) Read(key T) *H {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.collection == nil {
		return nil
	}

	if val, ok := b.collection[key]; ok {
		return &val
	}

	return nil
}

func (b *buffer[T, H]) ReadAll() []H {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.collection == nil {
		return nil
	}

	var values []H

	for _, val := range b.collection {
		values = append(values, val)
	}

	return values
}

func (b *buffer[T, H]) Write(key T, value H) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.collection == nil {
		b.collection = make(map[T]H)
	}

	b.collection[key] = value
}

func (b *buffer[T, H]) Remove(key T) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.collection == nil {
		b.collection = make(map[T]H)
	}

	delete(b.collection, key)
}

func (oic *OuterInstanceConfig) FromOuter() *OuterInstanceConfig {
	return utils.PtrOf(oic.Inner().Outer())
}

func (oic *OuterInstanceConfig) Inner() InstanceConfig {
	var _config = defaultInstanceConfig

	if oic == nil {
		return _config
	}

	if oic.maxOperationQueueSize != nil {
		_config.maxOperationQueueSize = *oic.maxOperationQueueSize
	}

	if oic.maxLogQueueSize != nil {
		_config.maxLogQueueSize = *oic.maxLogQueueSize
	}

	if oic.maxProductionChannelItems != nil {
		_config.maxProductionChannelItems = *oic.maxProductionChannelItems
	}

	if oic.maxStdoutWriterChannelItems != nil {
		_config.maxStdoutWriterChannelItems = *oic.maxStdoutWriterChannelItems
	}

	if oic.maxFileWriterChannelItems != nil {
		_config.maxFileWriterChannelItems = *oic.maxFileWriterChannelItems
	}

	if oic.maxExternalWriterChannelItems != nil {
		_config.maxExternalWriterChannelItems = *oic.maxExternalWriterChannelItems
	}

	if oic.maxStdoutWriters != nil {
		_config.maxStdoutWriters = *oic.maxStdoutWriters
	}

	if oic.maxFileWriters != nil {
		_config.maxFileWriters = *oic.maxFileWriters
	}

	if oic.maxExternalWriters != nil {
		_config.maxExternalWriters = *oic.maxExternalWriters
	}

	if oic.idleChannelFlushInterval != nil {
		_config.idleChannelFlushInterval = *oic.idleChannelFlushInterval
	}

	if oic.systemStatsCheckInterval != nil {
		_config.systemStatsCheckInterval = *oic.systemStatsCheckInterval
	}

	if oic.writeToStdout != nil {
		_config.writeToStdout = *oic.writeToStdout
	}

	if oic.writeToFile != nil {
		_config.writeToFile = *oic.writeToFile
	}

	if oic.writeToExternal != nil {
		_config.writeToExternal = *oic.writeToExternal
	}

	if oic.fileWriterBasePath != nil {
		_config.fileWriterBasePath = *oic.fileWriterBasePath
	}

	if oic.fileLogsDirectoryGranularity != nil {
		_config.fileLogsDirectoryGranularity = *oic.fileLogsDirectoryGranularity
	}

	if oic.fileLogsDirectoryFormatLayout != nil {
		_config.fileLogsDirectoryFormatLayout = *oic.fileLogsDirectoryFormatLayout
	}

	if oic.operationsFileName != nil {
		_config.operationsFileName = *oic.operationsFileName
	}

	if oic.logsFileName != nil {
		_config.logsFileName = *oic.logsFileName
	}

	return _config
}

func (ic InstanceConfig) Outer() OuterInstanceConfig {
	return OuterInstanceConfig{
		maxOperationQueueSize:         &ic.maxOperationQueueSize,
		maxLogQueueSize:               &ic.maxLogQueueSize,
		maxProductionChannelItems:     &ic.maxProductionChannelItems,
		maxStdoutWriterChannelItems:   &ic.maxStdoutWriterChannelItems,
		maxFileWriterChannelItems:     &ic.maxFileWriterChannelItems,
		maxExternalWriterChannelItems: &ic.maxExternalWriterChannelItems,
		maxStdoutWriters:              &ic.maxStdoutWriters,
		maxFileWriters:                &ic.maxFileWriters,
		maxExternalWriters:            &ic.maxExternalWriters,
		idleChannelFlushInterval:      &ic.idleChannelFlushInterval,
		systemStatsCheckInterval:      &ic.systemStatsCheckInterval,
		writeToStdout:                 &ic.writeToStdout,
		writeToFile:                   &ic.writeToFile,
		writeToExternal:               &ic.writeToExternal,
		fileWriterBasePath:            &ic.fileWriterBasePath,
		fileLogsDirectoryGranularity:  &ic.fileLogsDirectoryGranularity,
		fileLogsDirectoryFormatLayout: &ic.fileLogsDirectoryFormatLayout,
		operationsFileName:            &ic.operationsFileName,
		logsFileName:                  &ic.logsFileName,
	}
}

func (l Log) String(r *Roga) string {
	return r.stdoutLogFormatter.Format(l)
}

func (o Operation) String(r *Roga) string {
	return r.stdoutOperationFormatter.Format(o)
}
