package roga

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dullkingsman/go-pkg/roga/writable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapWriter struct {
	basePath    string
	fileFormat  string
	logFiles    []string
	loggers     map[string]*zap.Logger
	fileEncoder zapcore.Encoder
	mu          sync.RWMutex
}

func NewZapWriter(basePath string, fileFormat string, logFiles ...string) (*ZapWriter, error) {
	zw := &ZapWriter{
		basePath:    basePath,
		fileFormat:  fileFormat,
		logFiles:    logFiles,
		loggers:     make(map[string]*zap.Logger),
		fileEncoder: newFileEncoder(),
	}

	if err := zw.initLoggers(); err != nil {
		return nil, err
	}

	go zw.scheduleDailyRefresh()

	return zw, nil
}

func (z *ZapWriter) initLoggers() error {
	z.mu.Lock()
	defer z.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	dailyPath := filepath.Join(z.basePath, today)

	if err := os.MkdirAll(dailyPath, 0755); err != nil {
		return fmt.Errorf("could not create daily folder: %w", err)
	}

	for _, fileName := range z.logFiles {
		fullName := fileName + z.fileFormat
		logPath := filepath.Join(dailyPath, fullName)

		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("could not open log file %s: %w", fullName, err)
		}

		writer := zapcore.AddSync(file)

		core := zapcore.NewCore(
			z.fileEncoder,
			writer,
			zap.InfoLevel,
		)

		logger := zap.New(core)
		z.loggers[fullName] = logger
	}

	return nil
}

func (z *ZapWriter) WriteToFile(items ...writable.Writable) error {
	for _, item := range items {
		if err := z.ToFile(item); err != nil {
			return err
		}
	}
	return nil
}

func (z *ZapWriter) ToFile(item writable.Writable) error {
	switch v := item.(type) {
	case *Log:
		var logFile string
		switch v.Type {
		case TypeAudit:
			logFile = "audit"
		case TypeEvent:
			logFile = "events"
		default:
			logFile = "normal"
		}
		fullName := logFile + z.fileFormat

		z.mu.RLock()
		logger, ok := z.loggers[fullName]
		z.mu.RUnlock()
		if !ok {
			return fmt.Errorf("no logger configured for log type: %s", v.Type)
		}

		fields := []zap.Field{
			zap.String("type", "log"),
			zap.Any("log", v),
		}

		switch v.Level {
		case LevelFatal:
			logger.Fatal("", fields...)
		case LevelError:
			logger.Error("", fields...)
		case LevelWarn:
			logger.Warn("", fields...)
		case LevelInfo:
			logger.Info("", fields...)
		case LevelDebug:
			logger.Debug("", fields...)
		default:
			logger.Info("", fields...)
		}

	case *Operation:
		fullName := "operations" + z.fileFormat

		z.mu.RLock()
		logger, ok := z.loggers[fullName]
		z.mu.RUnlock()
		if !ok {
			return fmt.Errorf("no logger configured for: %s", fullName)
		}

		logger.Info("operation",
			zap.String("type", "operation"),
			zap.Any("operation", v),
		)

	default:
		return fmt.Errorf("unsupported type: %T", v)
	}

	return nil
}

func (z *ZapWriter) WriteToStdout(items ...writable.Writable) {
}

func newFileEncoder() zapcore.Encoder {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "timestamp"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encCfg)
}

func (z *ZapWriter) scheduleDailyRefresh() {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 1, 0, now.Location())
		time.Sleep(time.Until(next))

		_ = z.initLoggers()
	}
}
