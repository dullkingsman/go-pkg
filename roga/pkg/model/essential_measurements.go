package model

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
	"time"
)

type EssentialMeasurements struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func (c *EssentialMeasurements) Json() ([]byte, error) {
	if c == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	mj.WriteTimeField("startTime", c.StartTime, nil)
	mj.WriteTimeField("endTime", c.EndTime, nil)

	return mj.End(), nil
}
