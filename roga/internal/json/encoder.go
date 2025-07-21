package json

import (
	"bytes"
	"encoding/json"
	"github.com/dullkingsman/go-pkg/utils"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type ManualJsonEncoder struct {
	bytes.Buffer
	array bool
}

func NewManualJson(array ...bool) *ManualJsonEncoder {
	var mj = &ManualJsonEncoder{}

	if len(array) > 0 && array[0] {
		mj.array = true
		mj.StartArray()
	} else {
		mj.StartObject()
	}

	return mj
}

func (mj *ManualJsonEncoder) StartObject() {
	mj.WriteByte('{')
}

func (mj *ManualJsonEncoder) StartArray() {
	mj.WriteByte('[')
}

func (mj *ManualJsonEncoder) isOmittable(omittable ...bool) bool {
	var _omittable bool

	if len(omittable) > 0 {
		_omittable = omittable[0]
	}

	return _omittable
}

func (mj *ManualJsonEncoder) WriteField(field string) *ManualJsonEncoder {
	if mj.Len() > 1 {
		mj.WriteByte(',')
	}

	mj.WriteString(field)

	return mj
}

func (mj *ManualJsonEncoder) WriteStringField(name string, value string, omittable ...bool) *ManualJsonEncoder {
	if mj.isOmittable(omittable...) && value == "" {
		return nil
	}

	mj.WriteField(`"` + name + `":"` + strings.ReplaceAll(value, `"`, `\"`) + `"`)

	return mj
}

func (mj *ManualJsonEncoder) WriteUuidField(name string, value uuid.UUID, omittable ...bool) *ManualJsonEncoder {
	if mj.isOmittable(omittable...) && value == uuid.Nil {
		return nil
	}

	mj.WriteField(`"` + name + `":"` + value.String() + `"`)

	return mj
}

func (mj *ManualJsonEncoder) WriteTimeField(name string, value time.Time, format *string, omittable ...bool) *ManualJsonEncoder {
	if mj.isOmittable(omittable...) && value.IsZero() {
		return nil
	}

	var tmp = ""

	if format != nil {
		tmp = `"` + value.Format(*format) + `"`
	} else {
		tmp = strconv.FormatInt(value.UnixNano(), 10)
	}

	mj.WriteField(`"` + name + `":` + tmp)

	return mj
}

func (mj *ManualJsonEncoder) WriteInt64Field(name string, value int64, omittable ...bool) *ManualJsonEncoder {
	if mj.isOmittable(omittable...) && value == 0 {
		return nil
	}

	mj.WriteField(`"` + name + `":` + strconv.FormatInt(value, 10))

	return mj
}

func (mj *ManualJsonEncoder) WriteUint64Field(name string, value uint64, omittable ...bool) *ManualJsonEncoder {
	if mj.isOmittable(omittable...) && value == 0 {
		return nil
	}

	mj.WriteField(`"` + name + `":` + strconv.FormatUint(value, 10))

	return mj
}

func (mj *ManualJsonEncoder) WriteFloat64Field(name string, value float64, omittable ...bool) *ManualJsonEncoder {
	if mj.isOmittable(omittable...) && value == 0 {
		return nil
	}

	mj.WriteField(`"` + name + `":` + strconv.FormatFloat(value, 'f', -1, 64))

	return mj
}

func (mj *ManualJsonEncoder) WriteBooleanField(name string, value bool, omittable ...bool) *ManualJsonEncoder {
	if mj.isOmittable(omittable...) && value == false {
		return nil
	}

	mj.WriteField(`"` + name + `":` + strconv.FormatBool(value))

	return mj
}

func (mj *ManualJsonEncoder) WriteMarshalledJsonField(name string, value []byte, omittable ...bool) *ManualJsonEncoder {
	if mj.isOmittable(omittable...) && len(value) > 0 {
		return nil
	}

	if len(value) == 0 {
		return mj
	}

	mj.WriteField(`"` + name + `":` + string(value))

	return mj
}

func (mj *ManualJsonEncoder) WriteJsonField(name string, value interface{}, isArray bool, omittable ...bool) error {
	if mj.isOmittable(omittable...) && utils.UnderlyingValueIsNil(value) {
		return nil
	}

	if !utils.UnderlyingValueIsNil(value) {
		var _json, err = json.Marshal(value)

		if err != nil {
			utils.LogError("roga:encoder", "could not marshal data filed on log %v", err)
			return err
		}

		if len(_json) > 0 {
			mj.WriteField(`"` + name + `":` + string(_json))
		} else if isArray {
			mj.WriteField(`"` + name + `":[]`)
		} else {
			mj.WriteField(`"` + name + `":{}`)
		}
	}

	return nil
}

func (mj *ManualJsonEncoder) EndArray() *ManualJsonEncoder {
	mj.WriteByte(']')
	return mj
}

func (mj *ManualJsonEncoder) EndObject() *ManualJsonEncoder {
	mj.WriteByte('}')
	return mj
}

func (mj *ManualJsonEncoder) End() []byte {
	if mj.array {
		mj.EndArray()
	} else {
		mj.EndObject()
	}

	return mj.Bytes()
}
