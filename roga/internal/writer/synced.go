package writer

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

type SyncedWriter struct {
	bufio.Writer
	mu sync.Mutex
}

func NewSyncedWriter(writer io.Writer, size int) *SyncedWriter {
	return &SyncedWriter{
		Writer: *bufio.NewWriterSize(writer, size),
	}
}

func (sw *SyncedWriter) checkToFlush(in []byte) error {
	if len(in) > sw.Available() && sw.Buffered() > 0 {
		if err := sw.Flush(); err != nil {
			fmt.Println("failed to flush buffer:", err)
		}
	}

	//if sw.Writer.Buffered() > 0 {
	//	return sw.Flush()
	//}

	return nil
}

func (sw *SyncedWriter) Write(p []byte) (int, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if err := sw.checkToFlush(p); err != nil {
		return 0, err
	}

	return sw.Writer.Write(p)
}

func (sw *SyncedWriter) WriteString(s string) (int, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if err := sw.checkToFlush([]byte(s)); err != nil {
		return 0, err
	}

	return sw.Writer.WriteString(s)
}
