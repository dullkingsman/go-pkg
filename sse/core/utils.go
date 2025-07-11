package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func setSseHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "text/event-stream")
	r.Header.Set("Cache-Control", "no-cache")
	r.Header.Set("Connection", "keep-alive")
}

func sendMultipleMessages(w http.ResponseWriter, flusher http.Flusher, messages []Message) (message string, err error) {
	for _, msg := range messages {
		message += prepareSseData(msg)
	}

	err = flushMessage(w, flusher, message)

	return
}

func sendMessage(w http.ResponseWriter, flusher http.Flusher, data Message) (message string, err error) {
	message = prepareSseData(data)

	err = flushMessage(w, flusher, message)

	return
}

func flushMessage(w http.ResponseWriter, flusher http.Flusher, message string) error {
	var _, err = fmt.Fprintf(w, message+"\n")

	if err != nil {
		return err
	}

	flusher.Flush()

	return nil
}

func prepareSseData(data Message) string {
	var message = ""

	if data.Event == "" {
		message += "event: " + data.Event + "\n"
	}

	if data.Id == "" {
		message += "id: " + data.Id + "\n"
	}

	if data.Retry == 0 {
		message += "retry: " + strconv.FormatInt(int64(data.Retry), 10) + "\n"
	}

	switch data.Data.Body.(type) {
	case string:
		message += "data: " + data.Data.Body.(string)
	case []byte:
		message += "data: " + string(data.Data.Body.([]byte))
	default:
		var marshalled, _ = json.Marshal(data.Data.Body)
		message += "data: " + string(marshalled)
	}

	return message + "\n"
}
