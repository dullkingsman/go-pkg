package model

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
)

type ExternalSystem struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
}

func (c *ExternalSystem) Json() ([]byte, error) {
	if c == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	mj.WriteStringField("id", c.Id)
	mj.WriteStringField("name", c.Name)

	return mj.End(), nil
}
