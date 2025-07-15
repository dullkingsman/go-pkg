package roga

import (
	"github.com/dullkingsman/go-pkg/utils"
	"time"
)

func (r *Roga) monitorAndFlushIdleChannels() {
	r.consumptionSync.Add(1)
	defer r.consumptionSync.Done()

	utils.LogInfo("roga:startup", "monitoring idle channels...")

	for {
		select {
		case <-r.idleChannelMonitorControls.stop:
			return
		case <-r.idleChannelMonitorControls.pause:
			select {
			case <-r.idleChannelMonitorControls.stop:
				return
			case <-r.idleChannelMonitorControls.resume:
				continue
			}
		case <-r.idleChannelMonitorControls.resume:
			continue
		default:
			if r.lastWrite.Before(time.Now().Add(-r.config.idleChannelFlushInterval * time.Second)) {
				r.Flush()
			}

			time.Sleep(r.config.idleChannelFlushInterval * time.Second)
		}
	}
}

func (r *Roga) monitorAndUpdateSystemMetrics() {
	r.consumptionSync.Add(1)
	defer r.consumptionSync.Done()

	utils.LogInfo("roga:startup", "monitoring system metrics...")

	for {
		select {
		case <-r.metricMonitorControls.stop:
			return
		case <-r.metricMonitorControls.pause:
			select {
			case <-r.metricMonitorControls.stop:
				return
			case <-r.metricMonitorControls.resume:
				continue
			}
		case <-r.metricMonitorControls.resume:
			continue
		default:
			SetCurrentSystemMetrics(r)

			time.Sleep(r.config.systemStatsCheckInterval * time.Second)
		}
	}
}

func SetCurrentSystemMetrics(r *Roga) {
	r.metricsLock.Lock()
	defer r.metricsLock.Unlock()

	var (
		cpuUsage, cpuErr                   = r.monitor.GetCPUUsage()
		totalMemory, freeMemory, memoryErr = r.monitor.GetMemoryStats()
		totalSwap, freeSwap, swapErr       = r.monitor.GetSwapStats()
		totalDisk, freeDisk, diskErr       = r.monitor.GetDiskStats("/")
	)

	if cpuErr == nil {
		r.currentSystemMetrics.CpuUsage = cpuUsage
	}

	if memoryErr == nil {
		r.context.System.Memory = totalMemory
		r.currentSystemMetrics.AvailableMemory = freeMemory
	}

	if swapErr == nil {
		r.context.System.SwapSize = totalSwap
		r.currentSystemMetrics.AvailableSwap = freeSwap
	}

	if diskErr == nil {
		r.context.System.DiskSize = totalDisk
		r.currentSystemMetrics.AvailableDisk = freeDisk
	}
}

func (r *Roga) getOperationContext() Context {
	r.metricsLock.RLock()
	defer r.metricsLock.RUnlock()

	return r.context
}

type DefaultMonitor struct{ Monitor }

func (d DefaultMonitor) GetCPUUsage() (usage float64, err error) {
	return getCPUUsage()
}

func (d DefaultMonitor) GetMemoryStats() (total, free uint64, err error) {
	return getMemoryStats()
}

func (d DefaultMonitor) GetSwapStats() (total, free uint64, err error) {
	return getSwapStats()
}

func (d DefaultMonitor) GetDiskStats(path string) (total, free uint64, err error) {
	return getDiskStats(path)
}
