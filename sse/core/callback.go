package sse

import "net/http"

type DefaultCallback struct{ Callback }

func (d DefaultCallback) OnSend(err error, formattedMessage string, messages []Message) {

}

func (d DefaultCallback) OnConnection(w http.ResponseWriter, r *http.Request, externalId string, topics []string, flusher http.Flusher) error {
	return nil
}

func (d DefaultCallback) OnDisconnection(w http.ResponseWriter, r *http.Request, externalId string) error {
	return nil
}
