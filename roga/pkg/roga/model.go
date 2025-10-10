package roga

import (
	"github.com/dullkingsman/go-pkg/roga/internal/_map"
	"github.com/dullkingsman/go-pkg/roga/writable"
)

type (
	Channels struct {
		operational ChannelGroup
		flush       ActionChannelGroup
		stop        ActionChannelGroup
	}

	ChannelGroup struct {
		production chan writable.Writable
		writing    WriteStreamChannels
	}

	ActionChannelGroup struct {
		production chan bool
		writing    WriteStreamActionChannels
	}

	WriteStreamChannels       _map.NestedMap[WriteStream, EntryType, chan writable.Writable]
	WriteStreamActionChannels _map.NestedMap[WriteStream, EntryType, chan bool]
)
