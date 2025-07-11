package sse

import "net/http"

type Server struct {
	Instance
	clients     clientStore
	topics      topicStore
	externalIds externalIdMap
	callback    Callback
}

type topicStore map[string][]string // topic -> []clientId

type externalIdMap map[string]string // externalId -> clientId

type clientStore map[string]ServerClient

type clientMessageHistory []clientMessageHistoryItem

type clientMessageHistoryItem struct {
	messages          []Message
	formattedMessages string
	sent              bool
}

type ServerClient struct {
	id             string
	externalId     string
	channel        chan []Message
	messageHistory clientMessageHistory
}

type Message struct {
	Event string      `json:"event"`
	Data  MessageData `json:"data"`
	Id    string      `json:"id"`
	Retry int         `json:"retry"` // milliseconds before retrying the connection
}

type MessageData struct {
	Headers Headers     `json:"headers,omitempty"`
	Body    interface{} `json:"body"`
}

type Headers map[string]string

type Callback interface {
	OnSend(err error, formattedMessage string, messages []Message)
	OnConnection(w http.ResponseWriter, r *http.Request, externalId string, topics []string, flusher http.Flusher) error
	onDisconnection(w http.ResponseWriter, r *http.Request, externalId string) error
}

type Instance interface {
	AddConnection(w http.ResponseWriter, r *http.Request, callback Callback) error
	RemoveConnection(externalId string) error
}
