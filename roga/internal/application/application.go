package application

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
	"log"
	"os"
	"runtime"
)

type Application struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Version     string `json:"version"`
	Env         string `json:"env"`
	Node        int    `json:"node"`
	Instance    int    `json:"instance"`
	Lang        string `json:"lang"`
	LangVersion string `json:"langVersion"`
	ProcessId   int    `json:"processId"`
}

var DefaultApplication = Application{
	Lang:        "go",
	LangVersion: runtime.Version(),
	ProcessId:   os.Getpid(),
}

func (a *Application) SetApplicationSpecifics(
	code string,
	version string,
	env string,
	node *int,
	instance *int,
	name ...string,
) {
	var _name = code

	if len(name) > 0 && name[0] != "" {
		_name = name[0]
	}

	if code == "" {
		log.Fatalln("no code provided for the application")
	}

	if version == "" {
		log.Fatalln("no version provided for the application")
	}

	if env == "" {
		log.Fatalln("non env provided for the application")
	}

	if node != nil {
		a.Node = *node
	}

	if instance != nil {
		a.Instance = *instance
	}

	a.Name = _name
	a.Code = code
	a.Version = version
	a.Env = env
}

func (a *Application) Json() ([]byte, error) {
	if a == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	mj.WriteStringField("name", a.Name)
	mj.WriteStringField("code", a.Code)
	mj.WriteStringField("version", a.Version)
	mj.WriteStringField("env", a.Env)
	mj.WriteInt64Field("node", int64(a.Node))
	mj.WriteInt64Field("instance", int64(a.Instance))
	mj.WriteStringField("lang", a.Lang)
	mj.WriteStringField("langVersion", a.LangVersion)
	mj.WriteInt64Field("processId", int64(a.ProcessId))

	return mj.End(), nil
}
