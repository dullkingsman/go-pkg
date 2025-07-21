package model

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
	"github.com/dullkingsman/go-pkg/utils"
)

type Client struct {
	Id        string  `json:"id"`
	Ip        *string `json:"ip,omitempty"`
	UserAgent *string `json:"userAgent,omitempty"`
}

func (c *Client) Json() ([]byte, error) {
	if c == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	mj.WriteStringField("id", c.Id)
	mj.WriteStringField("ip", utils.ValueOr(c.Ip, ""), true)
	mj.WriteStringField("userAgent", utils.ValueOr(c.UserAgent, ""), true)

	return mj.End(), nil
}
