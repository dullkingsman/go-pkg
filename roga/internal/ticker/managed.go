package ticker

import (
	"sync"
	"time"
)

type ManagedTicker struct {
	interval time.Duration
	stopCh   chan struct{}
	pauseCh  chan bool
	wg       *sync.WaitGroup
	mu       sync.RWMutex
	paused   bool
}

func NewManagedTicker(interval time.Duration, fn func(), wg *sync.WaitGroup) *ManagedTicker {
	var t = &ManagedTicker{
		interval: interval,
		stopCh:   make(chan struct{}),
		pauseCh:  make(chan bool),
	}

	if wg != nil {
		wg.Add(1)
	} else {
		t.wg = &sync.WaitGroup{}
		t.wg.Add(1)
	}

	go func() {
		if wg != nil {
			defer wg.Done()
		} else {
			defer t.wg.Done()
		}

		var ticker = time.NewTicker(t.interval)
		defer ticker.Stop()

		for {
			select {
			case <-t.stopCh:
				return
			case pause := <-t.pauseCh:
				t.mu.Lock()
				t.paused = pause
				t.mu.Unlock()
			case <-ticker.C:
				t.mu.RLock()
				if !t.paused {
					fn()
				}
				t.mu.RUnlock()
			}
		}
	}()

	return t
}

func (t *ManagedTicker) Pause() {
	t.pauseCh <- true
}

func (t *ManagedTicker) Resume() {
	t.pauseCh <- false
}

func (t *ManagedTicker) Stop() {
	close(t.stopCh)

	if t.wg != nil {
		t.wg.Wait()
	}
}
