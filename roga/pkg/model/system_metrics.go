package model

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
)

type SystemMetrics struct {
	CpuUsage        float64 `json:"cpuUsage"`
	AvailableMemory uint64  `json:"availableMemory"`
	AvailableDisk   uint64  `json:"availableDisk"`
	AvailableSwap   uint64  `json:"availableSwap"`
}

func (c *SystemMetrics) Json() ([]byte, error) {
	if c == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	mj.WriteFloat64Field("cpuUsage", c.CpuUsage)
	mj.WriteUint64Field("availableMemory", uint64(c.CpuUsage))
	mj.WriteUint64Field("availableDisk", uint64(c.CpuUsage))
	mj.WriteUint64Field("availableSwap", uint64(c.CpuUsage))

	return mj.End(), nil
}
