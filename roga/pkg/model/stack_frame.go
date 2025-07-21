package model

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
)

type StackFrame struct {
	File       string `json:"file"`
	Function   string `json:"function"`
	LineNumber int    `json:"lineNumber"`
}

func (c *StackFrame) Json() ([]byte, error) {
	if c == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	mj.WriteStringField("name", c.File)
	mj.WriteStringField("serial", c.Function)
	mj.WriteInt64Field("uuid", int64(c.LineNumber))

	return mj.End(), nil
}
