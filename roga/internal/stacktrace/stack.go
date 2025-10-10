package stacktrace

import (
	model2 "github.com/dullkingsman/go-pkg/roga/pkg/model"
	"runtime"
)

// GetStackFrames
func GetStackFrames(framesToSkip int) []model2.StackFrame {
	var (
		stack = make([]model2.StackFrame, 0)
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

		stack = append(stack, model2.StackFrame{
			File:       frame.File,
			Function:   frame.Function,
			LineNumber: frame.Line,
		})
	}

	return stack
}
