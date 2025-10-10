package roga

import (
	"github.com/dullkingsman/go-pkg/roga/internal/sysinfo"
)

type DefaultMonitor struct{}

func (d DefaultMonitor) GetCPUUsage() (usage float64, err error) {
	return sysinfo.GetCPUUsage()
}

func (d DefaultMonitor) GetMemoryStats() (total, free uint64, err error) {
	return sysinfo.GetMemoryStats()
}

func (d DefaultMonitor) GetSwapStats() (total, free uint64, err error) {
	return sysinfo.GetSwapStats()
}

func (d DefaultMonitor) GetDiskStats(path string) (total, free uint64, err error) {
	return sysinfo.GetDiskStats(path)
}

type Monitor interface {
	GetCPUUsage() (usage float64, err error)
	GetMemoryStats() (total, free uint64, err error)
	GetSwapStats() (total, free uint64, err error)
	GetDiskStats(path string) (total, free uint64, err error)
}
