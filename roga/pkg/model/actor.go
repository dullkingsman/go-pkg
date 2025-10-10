package model

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
)

type Actor struct {
	Type           ActorType       `json:"type"`
	Client         *Client         `json:"client,omitempty"`
	User           *User           `json:"user,omitempty"`
	ExternalSystem *ExternalSystem `json:"externalSystem,omitempty"`
}

type (
	ActorType uint
)

const (
	ActorTypeSystem         ActorType = 0 // system
	ActorTypeUser           ActorType = 1 // user
	ActorTypeExternalSystem ActorType = 2 // external system
)

func (a *Actor) Json() ([]byte, error) {
	if a == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	mj.WriteInt64Field("type", int64(a.Type))

	var client, err = a.Client.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("client", client, true)

	user, err := a.User.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("user", user, true)

	externalSystem, err := a.ExternalSystem.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("externalSystem", externalSystem, true)

	return mj.End(), nil
}
