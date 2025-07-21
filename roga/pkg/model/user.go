package model

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
	"github.com/dullkingsman/go-pkg/utils"
)

type User struct {
	Identifier      string  `json:"identifier"` // anything specific that can identify the user. E.g. if the user is not yet created the phone number and if they are, the id.
	Id              *string `json:"id,omitempty,omitempty"`
	IdType          *string `json:"idType,omitempty,omitempty"`
	SessionId       *string `json:"sessionId,omitempty,omitempty"`
	SessionIdType   *string `json:"sessionIdType,omitempty,omitempty"`
	Role            *string `json:"role,omitempty,omitempty"`
	PermissionLevel *string `json:"permissionLevel,omitempty,omitempty"`
	Type            *string `json:"type,omitempty,omitempty"`
	PhoneNumber     *string `json:"phoneNumber,omitempty,omitempty"`
	Email           *string `json:"email,omitempty,omitempty"`
	Username        *string `json:"username,omitempty,omitempty"`
}

func (u *User) Json() ([]byte, error) {
	if u == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	mj.WriteStringField("identifier", u.Identifier)
	mj.WriteStringField("id", utils.ValueOr(u.Id, ""), true)
	mj.WriteStringField("idType", utils.ValueOr(u.IdType, ""), true)
	mj.WriteStringField("sessionId", utils.ValueOr(u.SessionId, ""), true)
	mj.WriteStringField("sessionIdType", utils.ValueOr(u.SessionIdType, ""), true)
	mj.WriteStringField("role", utils.ValueOr(u.Role, ""), true)
	mj.WriteStringField("permissionLevel", utils.ValueOr(u.PermissionLevel, ""), true)
	mj.WriteStringField("type", utils.ValueOr(u.Type, ""), true)
	mj.WriteStringField("phoneNumber", utils.ValueOr(u.PhoneNumber, ""), true)
	mj.WriteStringField("email", utils.ValueOr(u.Email, ""), true)
	mj.WriteStringField("username", utils.ValueOr(u.Username, ""), true)

	return mj.End(), nil
}
