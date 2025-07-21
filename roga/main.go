package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	//roga "github.com/dullkingsman/go-pkg/roga/core"
	"github.com/dullkingsman/go-pkg/utils"
	"os"
	//"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var startedAllocation = time.Now()

	var (
		items = make([]interface{}, 1_000_000)
		//entryType = roga.EntryTypeEvent
		//parentId = uuid.New()
	)

	var allocationTook = time.Since(startedAllocation)

	var startedObjectCreation = time.Now()

	for i := 0; i < 1_000_000; i++ {
		items[i] = map[string]interface{}{
			"index":    i,
			"category": "performance_test",
			"success":  true,
			"message":  "This is a log message",
			//"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
			//"level": "info",
			//"id": uuid.New(),
			//"operation_id": parentId,
			//"actor": map[string]interface{}{
			//	"type": "system",
			//},
			//"event": "SomethingHappened",
			//"outcome": "Succeeded",
		}

		//log.Id = uuid.New()
		//log.OperationId = parentId

		//items[i] = log
	}

	var objectCreationTook = time.Since(startedObjectCreation)

	var startedWriting = time.Now()

	var bufIo = bufio.NewWriter(os.Stdout)

	bufio.NewWriterSize(bufIo, 256*1024)

	for _, item := range items {
		if item == nil {
			continue
		}

		var bs, _ = json.Marshal(item)

		//switch entryType {
		//case roga.EntryTypeOperation:
		//	operation := utils.SafeCastValue[roga.Operation](item)
		//	entry := utils.CyanString("op", true) + "(" + utils.GreyString(operation.Id.String(), true) + ") " + operation.String(roga.DefaultStdoutFormatter{})
		//
		//	bs = strings.TrimSpace(utils.FormatInfoLog(operation.Name, entry, true))
		//
		//case roga.EntryTypeAudit, roga.EntryTypeEvent, roga.EntryTypeLog:
		//	logItem := utils.SafeCastValue[roga.Log](item)
		//	var fmtFunc = utils.FormatInfoLog
		//
		//	switch logItem.Level {
		//	case roga.LevelFatal, roga.LevelError:
		//		fmtFunc = utils.FormatErrorLog
		//	case roga.LevelWarn:
		//		fmtFunc = utils.FormatWarnLog
		//	case roga.LevelInfo:
		//		fmtFunc = utils.FormatInfoLog
		//	case roga.LevelDebug:
		//		fmtFunc = utils.FormatDebugLog
		//	default:
		//		fmtFunc = utils.FormatInfoLog
		//	}
		//
		//	bs = strings.TrimSpace(fmtFunc(roga.EntryTypeName[entryType], logItem.String(roga.DefaultStdoutFormatter{}), true))
		//}

		if len(bs) > bufIo.Available() && bufIo.Buffered() > 0 {
			if err := bufIo.Flush(); err != nil {
				utils.LogError("main", "Failed to flush buffer: %v", err)
			}
		}

		if _, err := bufIo.Write(append(bs, '\n')); err != nil {
			utils.LogError("main", "Failed to write to buffer: %v", err)
		}
	}

	if err := bufIo.Flush(); err != nil {
		utils.LogError("main", "Failed to flush buffer: %v", err)
	}

	var writeTook = time.Since(startedWriting)

	var totalTook = allocationTook + objectCreationTook + writeTook

	fmt.Printf("\n")
	fmt.Printf("-----------------------------------------------------\n")
	fmt.Printf("Allocation took: %v\n", allocationTook)
	fmt.Printf("Object creation took: %v\n", objectCreationTook)
	fmt.Printf("Write took: %v\n", writeTook)
	fmt.Printf("Total took: %v\n", totalTook)
	fmt.Printf("-----------------------------------------------------\n")

	return

	// Configure Zap for high performance to stdout (JSON format)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Customize time format
	encoderConfig.TimeKey = "timestamp"                   // Use "timestamp" instead of "ts"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout), // Direct to stdout
		zapcore.InfoLevel,          // Log level
	)
	logger := zap.New(core)
	defer logger.Sync() // Flushes buffer, if any

	numLines := 1000000 // 1 million lines
	// numLines := 10000000 // 10 million lines

	start := time.Now()

	for i := 0; i < numLines; i++ {
		logger.Info("This is a log message",
			zap.Int("index", i),
			zap.String("category", "performance_test"),
			zap.Bool("success", true),
		)
	}

	elapsed := time.Since(start)
	fmt.Printf("Logged %d lines in %s (%.2f lines/second)\n", numLines, elapsed, float64(numLines)/elapsed.Seconds())
}
