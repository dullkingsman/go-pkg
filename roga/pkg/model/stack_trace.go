package model

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
)

type StackTrace struct {
	Crashed bool         `json:"crashed"`
	Frames  []StackFrame `json:"frames,omitempty"`
}

func (c *StackTrace) Json() ([]byte, error) {
	if c == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	mj.WriteBooleanField("crashed", c.Crashed)

	if len(c.Frames) > 0 {
		var mj2 = json.NewManualJson(true)

		for _, frame := range c.Frames {
			var _frame, err = frame.Json()

			if err != nil {
				return nil, err
			}

			mj2.Write(_frame)
		}

		mj.WriteMarshalledJsonField("frames", mj2.End(), true)
	}

	return mj.End(), nil
}
