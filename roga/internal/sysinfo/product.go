package sysinfo

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
	"github.com/dullkingsman/go-pkg/utils"
)

type ProductIdentifier struct {
	Name   *string `json:"name,omitempty"`
	Serial *string `json:"serial,omitempty"`
	Uuid   *string `json:"uuid,omitempty"`
}

func (c *ProductIdentifier) Json() ([]byte, error) {
	if c == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	mj.WriteStringField("name", utils.ValueOr(c.Name, ""), true)
	mj.WriteStringField("serial", utils.ValueOr(c.Serial, ""), true)
	mj.WriteStringField("uuid", utils.ValueOr(c.Uuid, ""), true)

	return mj.End(), nil
}
