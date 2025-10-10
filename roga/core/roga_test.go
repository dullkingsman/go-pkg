package roga

import (
	"github.com/dullkingsman/go-pkg/roga/internal"
	"github.com/dullkingsman/go-pkg/roga/pkg/model"
	"github.com/dullkingsman/go-pkg/roga/pkg/roga"
	"github.com/dullkingsman/go-pkg/roga/writable"
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/google/uuid"
	"os"
	"sync"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	TestRogaInitialization(t)
	TestBasicLogging(t)
	TestOperationManagement(t)
	TestFileWriting(t)
}

// TestRogaInitialization tests that Roga initializes correctly
func TestRogaInitialization(t *testing.T) {
	// Initialize Roga with default configuration
	r := roga.Init(roga.Config{
		Code:    "test",
		Version: "1.0.0",
		Env:     "test",
	})

	// Verify that Roga is initialized by checking if Start() doesn't panic
	//r.Start()

	// Clean up
	r.Stop()
}

func TestQueueConsumer(t *testing.T) {
	t.Run("stopes", func(t *testing.T) {
		var (
			size       = 1
			queue      = make(chan writable.Writable, size)
			stopChan   = make(chan bool)
			flushChan  = make(chan bool)
			dependency = make(chan writable.Writable, size)
			//dependentStop  = make(chan bool)
			//dependentFlush = make(chan bool)
			wg = &sync.WaitGroup{}
		)

		wg.Add(1)

		go internal.ConsumeQueue(
			"test_queue",
			queue,
			stopChan,
			flushChan,
			[]<-chan writable.Writable{dependency},
			[]chan<- bool{},
			[]chan<- bool{},
			wg,
			func(items []writable.Writable) {
				for _, item := range items {
					println(item.String(roga.DefaultStdoutFormatter{}))
				}
			},
		)

		queue <- roga.Log{
			Message: "test message",
		}

		flushChan <- true

		stopChan <- true

		wg.Wait()
	})
}

// TestBasicLogging tests basic logging functionality
func TestBasicLogging(t *testing.T) {
	// Initialize Roga
	r := roga.Init()
	r.Start()

	// Create a simple actor for our logs
	actor := model.Actor{
		Type: model.ActorTypeSystem,
	}

	// Create log arguments
	logArgs := roga.LogArgs{
		Message: "Test log message",
		Actor:   &actor,
	}

	// Test logging at different levels
	// We're just testing that these don't panic
	r.LogInfo(logArgs)
	r.LogWarn(logArgs)
	r.LogError(logArgs)
	r.LogDebug(logArgs)
	// We don't test Fatal as it would exit the program

	// Clean up
	r.Stop()
}

// TestBasicLogging tests basic logging functionality
func TestBasicEventCapture(t *testing.T) {
	// Initialize Roga
	r := roga.Init()
	r.Start()

	// time
	r.CaptureEvent(roga.EventLogArgs{
		roga.LogArgs{
			Event:   utils.PtrOf("SomethingHappened"),
			Outcome: utils.PtrOf("Succeeded"),
			Message: "Successful happenstance",
		},
	})

	r.LogInfo(roga.LogArgs{
		Message: "Test in root log message",
	})

	//time.Sleep(6 * time.Second)

	// Clean up
	r.Stop()
}

// TestPerformanceLogging tests basic logging functionality
func TestPerformanceLogging(t *testing.T) {
	// Initialize Roga
	r := roga.Init(roga.Config{
		Code:    "test",
		Version: "1.0.0",
		Env:     "test",
		InstanceConfig: &roga.OuterInstanceConfig{
			WriteToFile: utils.PtrOf(false),
		},
	})
	r.Start()

	for i := 0; i < 1_000_000; i++ {
		r.CaptureEvent(roga.EventLogArgs{
			roga.LogArgs{
				Event:   utils.PtrOf("SomethingHappened"),
				Outcome: utils.PtrOf("Succeeded"),
				Message: "Successful happenstance",
			},
		})
	}

	//var startedAllocation = time.Now()
	//
	//var (
	//	items     = make([]Writable, 1_000_000)
	//	entryType = EntryTypeEvent
	//)
	//
	//var allocationTook = time.Since(startedAllocation)
	//
	//var startedObjectCreation = time.Now()
	//
	//for i := 0; i < 1_000_000; i++ {
	//	items[i] = EventLogArgs{
	//		LogArgs{
	//			Event:   utils.PtrOf("SomethingHappened"),
	//			Outcome: utils.PtrOf("Succeeded"),
	//			Message: "Successful happenstance",
	//		},
	//	}.ToLog()
	//}
	//
	//var objectCreationTook = time.Since(startedObjectCreation)
	//
	//var startedWriting = time.Now()
	//
	//for _, item := range items {
	//	if item == nil {
	//		continue
	//	}
	//
	//	switch entryType {
	//	case EntryTypeOperation:
	//		operation := utils.SafeCastValue[Operation](item)
	//		entry := utils.CyanString("op") + "(" + utils.GreyString(operation.Id.String()) + ") " + operation.String(DefaultStdoutFormatter{})
	//		fmt.Printf(strings.TrimSpace(operation.Name + " " + entry + "\n"))
	//
	//	case EntryTypeAudit, EntryTypeEvent, EntryTypeLog:
	//		logItem := utils.SafeCastValue[Log](item)
	//		//var fmtFunc = utils.FormatInfoLog
	//
	//		//switch logItem.Level {
	//		//case LevelFatal, LevelError:
	//		//	fmtFunc = utils.FormatErrorLog
	//		//case LevelWarn:
	//		//	fmtFunc = utils.FormatWarnLog
	//		//case LevelInfo:
	//		//	fmtFunc = utils.FormatInfoLog
	//		//case LevelDebug:
	//		//	fmtFunc = utils.FormatDebugLog
	//		//default:
	//		//	fmtFunc = utils.FormatInfoLog
	//		//}
	//
	//		fmt.Printf(strings.TrimSpace(EntryTypeName[entryType] + " " + logItem.String(DefaultStdoutFormatter{})))
	//	}
	//
	//}
	//
	//var writeTook = time.Since(startedWriting)
	//
	//var totalTook = allocationTook + objectCreationTook + writeTook
	//
	//fmt.Printf("-----------------------------------------------------\n")
	//fmt.Printf("-----------------------------------------------------\n")
	//fmt.Printf("Allocation took: %v\n", allocationTook)
	//fmt.Printf("Object creation took: %v\n", objectCreationTook)
	//fmt.Printf("Write took: %v\n", writeTook)
	//fmt.Printf("Total took: %v\n", totalTook)
	//fmt.Printf("-----------------------------------------------------\n")
	//fmt.Printf("Average allocation took: %v\n", allocationTook.Seconds()/float64(1_000_000))
	//fmt.Printf("Average object creation took: %v\n", objectCreationTook.Seconds()/float64(1_000_000))
	//fmt.Printf("Average write took: %v\n", writeTook.Seconds()/float64(1_000_000))
	//fmt.Printf("Average total took: %v\n", totalTook.Seconds()/float64(1_000_000))
	//fmt.Printf("-----------------------------------------------------\n")
	//fmt.Printf("-----------------------------------------------------\n")

	// Clean up
	r.Stop()
}

// TestBasicOperationLogging tests basic logging functionality
func TestBasicOperationLogging(t *testing.T) {
	// Create a temporary directory for logs
	var tempDir = "./roga_test/basic_operation_logging/"
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	//defer os.RemoveAll(tempDir)

	// Configure Roga to write to files
	customConfig := roga.Config{
		InstanceConfig: &roga.OuterInstanceConfig{
			WriteToStdout:      utils.PtrOf(false),
			WriteToFile:        utils.PtrOf(true),
			WriteToExternal:    utils.PtrOf(false),
			FileWriterBasePath: &tempDir,
		},
	}

	r := roga.Init(customConfig)
	r.Start()

	r.LogInfo(roga.LogArgs{
		Message: "Test in root log message",
		Actor: &model.Actor{
			Type: model.ActorTypeSystem,
		},
	})

	op := r.BeginOperation(roga.OperationArgs{
		Name: "TestOperation",
	})

	op.LogInfo(roga.LogArgs{
		Message: "Test in nested op log message",
		Actor: &model.Actor{
			Type: model.ActorTypeSystem,
		},
	})

	var op2 = op.BeginOperation(roga.OperationArgs{
		Name: "TestOperation2",
		Actor: &model.Actor{
			Type: model.ActorTypeUser,
			User: &model.User{
				Identifier:      uuid.New().String(),
				Id:              utils.PtrOf(uuid.New().String()),
				IdType:          utils.PtrOf("UUID"),
				SessionId:       utils.PtrOf(uuid.New().String()),
				SessionIdType:   utils.PtrOf("UUID"),
				Role:            nil,
				PermissionLevel: utils.PtrOf("BASIC"),
				Type:            utils.PtrOf("CUSTOMER"),
				PhoneNumber:     utils.PtrOf("+251960621337"),
				Email:           utils.PtrOf("daniel@gebta.app"),
			},
		},
	})

	op2.AuditAction(roga.AuditLogArgs{
		roga.LogArgs{
			Event:   utils.PtrOf("Login"),
			Outcome: utils.PtrOf("Succeeded"),
			Message: "Successful login",
		},
	})

	op.CaptureEvent(roga.EventLogArgs{
		roga.LogArgs{
			Event:   utils.PtrOf("SomethingHappened"),
			Outcome: utils.PtrOf("Succeeded"),
		},
	})

	if op == nil {
		t.Fatal("Expected operation to be created, got nil")
	}

	if op.Name != "TestOperation" {
		t.Errorf("Expected operation name to be TestOperation, got %s", op.Name)
	}

	op2.EndOperation()

	op.EndOperation()

	time.Sleep(2 * time.Second)

	r.Stop()

	//// Verify operation was logged
	//var logsBaseDir = tempDir + time2.GetTimeRoundedTo(
	//	r.config.fileLogsDirectoryGranularity,
	//).UTC().Format(
	//	r.config.fileLogsDirectoryFormatLayout,
	//)

	//files, err := os.ReadDir(logsBaseDir)
	//if err != nil {
	//	t.Fatalf("Failed to read temp directory: %v", err)
	//}

	//foundOperationsFile := false
	//for _, file := range files {
	//	if file.Name() == roga.DefaultOperationsFileName {
	//		foundOperationsFile = true
	//		break
	//	}
	//}

	//if !foundOperationsFile {
	//	t.Errorf("Expected to find operations file %s", roga.DefaultOperationsFileName)
	//}
}

// TestOperationManagement tests the creation and management of operations
func TestOperationManagement(t *testing.T) {
	// Initialize Roga
	r := roga.Init()
	r.Start()

	// Create a simple actor for our operations
	actor := model.Actor{
		Type: model.ActorTypeSystem,
	}

	// Create operation arguments
	opArgs := roga.OperationArgs{
		Name:  "TestOperation",
		Actor: &actor,
	}

	// Begin an operation
	op := r.BeginOperation(opArgs)
	if op == nil {
		t.Fatal("Expected operation to be created, got nil")
	}
	if op.Name != "TestOperation" {
		t.Errorf("Expected operation name to be %s, got %s", "TestOperation", op.Name)
	}

	// ProduceLog within the operation
	logArgs := roga.LogArgs{
		Message: "Test log within operation",
		Actor:   &actor,
	}
	op.LogInfo(logArgs)

	// Create a nested operation
	nestedOpArgs := roga.OperationArgs{
		Name:  "NestedOperation",
		Actor: &actor,
	}
	nestedOp := op.BeginOperation(nestedOpArgs)
	if nestedOp == nil {
		t.Fatal("Expected nested operation to be created, got nil")
	}
	if nestedOp.Name != "NestedOperation" {
		t.Errorf("Expected nested operation name to be %s, got %s", "NestedOperation", nestedOp.Name)
	}
	if nestedOp.ParentId == nil {
		t.Errorf("Expected nested operation to have a parent ID")
	}

	// End the nested operation
	nestedOp.EndOperation()

	// End the parent operation
	op.EndOperation()

	// Clean up
	r.Stop()
}

// TestFileWriting tests writing logs and operations to files
func TestFileWriting(t *testing.T) {
	// Create a temporary directory for logs
	tempDir, err := os.MkdirTemp("", "roga_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a custom configuration that writes to files
	customConfig := roga.Config{
		InstanceConfig: &roga.OuterInstanceConfig{
			WriteToStdout:      utils.PtrOf(false),
			WriteToFile:        utils.PtrOf(true),
			WriteToExternal:    utils.PtrOf(false),
			FileWriterBasePath: &tempDir,
		},
	}

	// Initialize Roga with custom configuration
	r := roga.Init(customConfig)
	r.Start()

	// Create a simple actor for our logs and operations
	actor := model.Actor{
		Type: model.ActorTypeSystem,
	}

	// Create and log some test data
	logArgs := roga.LogArgs{
		Message: "Test file writing",
		Actor:   &actor,
	}
	r.LogInfo(logArgs)

	opArgs := roga.OperationArgs{
		Name:  "TestFileOperation",
		Actor: &actor,
	}
	op := r.BeginOperation(opArgs)
	op.EndOperation()

	// Stop Roga
	r.Stop()

	time.Sleep(2 * time.Second)

	//var logsBaseDir = tempDir + internal.getCurrentTimeRoundedTo(
	//	r.config.fileLogsDirectoryGranularity,
	//).UTC().Format(
	//	r.config.fileLogsDirectoryFormatLayout,
	//)
	//
	//// Verify that log files were created
	//files, err := os.ReadDir(logsBaseDir)
	//if err != nil {
	//	t.Fatalf("Failed to read temp directory: %v", err)
	//}
	//
	//if len(files) == 0 {
	//	t.Error("Expected log files to be created, but directory is empty")
	//}
	//
	//foundOperationsFile := false
	//foundLogsFile := false
	//
	//for _, file := range files {
	//	if file.Name() == roga.DefaultOperationsFileName {
	//		foundOperationsFile = true
	//	}
	//	if file.Name() == "normal."+roga.DefaultLogsFileName {
	//		foundLogsFile = true
	//	}
	//}
	//
	//if !foundLogsFile {
	//	t.Errorf("Expected to find logs file %s", "normal."+roga.DefaultLogsFileName)
	//}
	//
	//if !foundOperationsFile {
	//	t.Errorf("Expected to find operations file %s", roga.DefaultOperationsFileName)
	//}
}
