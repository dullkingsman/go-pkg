package core

import "time"

func (r *Roga) monitorAndFlushIdleChannels() {
	r.consumptionSync.Add(1)
	defer r.consumptionSync.Done()

	for {
		select {
		case <-r.idleChannelMonitorControls.stop:
			return
		case <-r.idleChannelMonitorControls.pause:
			<-r.idleChannelMonitorControls.resume
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

	for {
		select {
		case <-r.metricMonitorControls.stop:
			return
		case <-r.metricMonitorControls.pause:
			<-r.metricMonitorControls.resume
		default:
			var (
				cpuUsage, cpuErr                   = r.monitor.GetCPUUsage()
				totalMemory, freeMemory, memoryErr = r.monitor.GetMemoryStats()
				totalSwap, freeSwap, swapErr       = r.monitor.GetSwapStats()
				totalDisk, freeDisk, diskErr       = r.monitor.GetDiskStats("/")
			)

			r.metricsLock.Lock()

			if cpuErr == nil {
				r.context.Environment.SystemEnvironment.CpuUsage = cpuUsage
			}

			if memoryErr == nil {
				r.context.SystemSpecifications.Memory = totalMemory
				r.context.Environment.SystemEnvironment.AvailableMemory = freeMemory
			}

			if swapErr == nil {
				r.context.SystemSpecifications.SwapSize = totalSwap
				r.context.Environment.SystemEnvironment.AvailableSwap = freeSwap
			}

			if diskErr == nil {
				r.context.SystemSpecifications.DiskSize = totalDisk
				r.context.Environment.SystemEnvironment.AvailableDisk = freeDisk
			}

			r.metricsLock.Unlock()

			time.Sleep(r.config.systemStatsCheckInterval * time.Second)
		}
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
