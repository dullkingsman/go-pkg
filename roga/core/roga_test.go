package roga

import (
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/google/uuid"
	"os"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	TestRogaInitialization(t)
	TestBasicLogging(t)
	TestOperationManagement(t)
	TestMonitoring(t)
	TestFileWriting(t)
}

// TestRogaInitialization tests that Roga initializes correctly
func TestRogaInitialization(t *testing.T) {
	// Initialize Roga with default configuration
	r := Init()

	utils.LogInfo("test:roga", "Initialized Roga")

	// Verify that Roga is initialized by checking if Start() doesn't panic
	r.Start()

	utils.LogInfo("test:roga", "Started Roga")

	// Clean up
	r.Stop()

	utils.LogInfo("test:roga", "Stopped Roga")
	utils.LogInfo("test:roga", "Finished Roga Initialization")
}

// TestBasicLogging tests basic logging functionality
func TestBasicLogging(t *testing.T) {
	// Initialize Roga
	r := Init()
	r.Start()

	// Create a simple actor for our logs
	actor := Actor{
		Type: ActorTypeSystem,
	}

	// Create log arguments
	logArgs := LogArgs{
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
	r.Stop(true)
}

// TestPerformanceLogging tests basic logging functionality
func TestPerformanceLogging(t *testing.T) {
	// Initialize Roga
	r := Init()
	r.Start()

	// Create a simple actor for our logs
	actor := Actor{
		Type: ActorTypeSystem,
	}

	// Create log arguments
	logArgs := LogArgs{
		Message: "Test log message",
		Actor:   &actor,
	}

	// Test logging at different levels
	// We're just testing that these don't panic
	for i := 0; i < 100_000; i++ {
		r.LogInfo(logArgs)
	}
	// We don't test Fatal as it would exit the program

	//time.Sleep(10 * time.Second)

	// Clean up
	r.Stop(true)
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
	customConfig := Config{
		Instance: &OuterInstanceConfig{
			writeToStdout:      utils.PtrOf(false),
			writeToFile:        utils.PtrOf(true),
			writeToExternal:    utils.PtrOf(false),
			fileWriterBasePath: &tempDir,
		},
	}

	r := Init(customConfig)
	r.Start()

	r.LogInfo(LogArgs{
		Message: "Test in root log message",
		Actor: &Actor{
			Type: ActorTypeSystem,
		},
	})

	op := r.BeginOperation(OperationArgs{
		Name: "TestOperation",
	})

	op.LogInfo(LogArgs{
		Message: "Test in nested op log message",
		Actor: &Actor{
			Type: ActorTypeSystem,
		},
	})

	var op2 = r.BeginOperation(OperationArgs{
		Name: "TestOperation2",
		Actor: &Actor{
			Type: ActorTypeUser,
			User: &User{
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

	op2.AuditAction(AuditLogArgs{
		LogArgs{
			Event:   utils.PtrOf("Login"),
			Outcome: utils.PtrOf("Succeeded"),
			Message: "Successful login",
		},
	})

	op.CaptureEvent(EventLogArgs{
		LogArgs{
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

	op.EndOperation()

	time.Sleep(2 * time.Second)

	r.Stop(true)

	// Verify operation was logged
	var logsBaseDir = tempDir + getCurrentTimeRoundedTo(
		r.config.fileLogsDirectoryGranularity,
	).UTC().Format(
		r.config.fileLogsDirectoryFormatLayout,
	)

	files, err := os.ReadDir(logsBaseDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}

	foundOperationsFile := false
	for _, file := range files {
		if file.Name() == DefaultOperationsFileName {
			foundOperationsFile = true
			break
		}
	}

	if !foundOperationsFile {
		t.Errorf("Expected to find operations file %s", DefaultOperationsFileName)
	}
}

// TestOperationManagement tests the creation and management of operations
func TestOperationManagement(t *testing.T) {
	// Initialize Roga
	r := Init()
	r.Start()

	// Create a simple actor for our operations
	actor := Actor{
		Type: ActorTypeSystem,
	}

	// Create operation arguments
	opArgs := OperationArgs{
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

	// Log within the operation
	logArgs := LogArgs{
		Message: "Test log within operation",
		Actor:   &actor,
	}
	op.LogInfo(logArgs)

	// Create a nested operation
	nestedOpArgs := OperationArgs{
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
	r.Stop(true)
}

// TestMonitoring tests the system monitoring functionality
func TestMonitoring(t *testing.T) {
	// Create a custom configuration
	customConfig := Config{
		Instance: &OuterInstanceConfig{
			systemStatsCheckInterval: utils.PtrOf(time.Duration(1)), // Check every second for faster testing
		},
	}

	// Initialize Roga with custom configuration
	r := Init(customConfig)
	r.Start()

	// Wait for metrics to be collected
	time.Sleep(2 * time.Second)

	// Test pausing and resuming monitoring
	r.PauseSystemMonitoring()
	time.Sleep(1 * time.Second)
	r.ResumeSystemMonitoring()
	time.Sleep(1 * time.Second)

	// Stop monitoring
	r.StopSystemMonitoring()

	// Clean up
	r.Stop(true)
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
	customConfig := Config{
		Instance: &OuterInstanceConfig{
			writeToStdout:      utils.PtrOf(false),
			writeToFile:        utils.PtrOf(true),
			writeToExternal:    utils.PtrOf(false),
			fileWriterBasePath: &tempDir,
		},
	}

	// Initialize Roga with custom configuration
	r := Init(customConfig)
	r.Start()

	// Create a simple actor for our logs and operations
	actor := Actor{
		Type: ActorTypeSystem,
	}

	// Create and log some test data
	logArgs := LogArgs{
		Message: "Test file writing",
		Actor:   &actor,
	}
	r.LogInfo(logArgs)

	opArgs := OperationArgs{
		Name:  "TestFileOperation",
		Actor: &actor,
	}
	op := r.BeginOperation(opArgs)
	op.EndOperation()

	// Stop Roga
	r.Stop(true)

	time.Sleep(2 * time.Second)

	var logsBaseDir = tempDir + getCurrentTimeRoundedTo(
		r.config.fileLogsDirectoryGranularity,
	).UTC().Format(
		r.config.fileLogsDirectoryFormatLayout,
	)

	// Verify that log files were created
	files, err := os.ReadDir(logsBaseDir)
	if err != nil {
		t.Fatalf("Failed to read temp directory: %v", err)
	}

	if len(files) == 0 {
		t.Error("Expected log files to be created, but directory is empty")
	}

	foundOperationsFile := false
	foundLogsFile := false

	for _, file := range files {
		if file.Name() == DefaultOperationsFileName {
			foundOperationsFile = true
		}
		if file.Name() == "normal."+DefaultLogsFileName {
			foundLogsFile = true
		}
	}

	if !foundLogsFile {
		t.Errorf("Expected to find logs file %s", "normal."+DefaultLogsFileName)
	}

	if !foundOperationsFile {
		t.Errorf("Expected to find operations file %s", DefaultOperationsFileName)
	}
}
