package model

import (
	"github.com/dullkingsman/go-pkg/roga/internal/application"
	"github.com/dullkingsman/go-pkg/roga/internal/json"
	"github.com/dullkingsman/go-pkg/roga/internal/sysinfo"
)

type Context struct {
	Application application.Application      `json:"application"`
	System      sysinfo.SystemSpecifications `json:"system"`
}

var DefaultContext = Context{
	Application: application.DefaultApplication,
	System:      sysinfo.GetSystemSpecifications(),
}

func (c *Context) Json() ([]byte, error) {
	if c == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	var application, err = c.Application.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("application", application, true)

	system, err := c.System.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("system", system, true)

	return mj.End(), nil
}
