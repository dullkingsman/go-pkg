package roga

import (
	"github.com/dullkingsman/go-pkg/roga/writable"
	"github.com/dullkingsman/go-pkg/utils"
)

type Config struct {
	Name            string
	Code            string
	Version         string
	Env             string
	Node            *int
	Instance        *int
	InstanceConfig  *OuterInstanceConfig
	StdoutFormatter writable.Formatter
	Producer        Producer
	Monitor         Monitor
	Writer          Writer
}

var defaultConfig = Config{
	InstanceConfig:  utils.PtrOf(DefaultInstanceConfig.Outer()),
	StdoutFormatter: &DefaultStdoutFormatter{},
	Producer:        &DefaultProducer{},
	Monitor:         &DefaultMonitor{},
	Writer:          &DefaultWriter{},
}

func (c Config) FromIncoming(incomingConfig ...Config) Config {
	var _config = c

	if len(incomingConfig) == 0 {
		return _config
	}

	var cfg = incomingConfig[0]

	_config.Name = cfg.Name
	_config.Code = cfg.Code
	_config.Version = cfg.Version
	_config.Env = cfg.Env
	_config.Node = cfg.Node
	_config.Instance = cfg.Instance

	if cfg.InstanceConfig != nil {
		_config.InstanceConfig = utils.PtrOf(cfg.InstanceConfig.Inner().Outer())
	}

	if cfg.Producer != nil {
		_config.Producer = cfg.Producer
	}

	if cfg.Monitor != nil {
		_config.Monitor = cfg.Monitor
	}

	if cfg.Writer != nil {
		_config.Writer = cfg.Writer
	}

	if cfg.StdoutFormatter != nil {
		_config.StdoutFormatter = cfg.StdoutFormatter
	}

	return _config
}
