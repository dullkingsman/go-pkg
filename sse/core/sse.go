package sse

import (
	"errors"
	"github.com/dullkingsman/go-pkg/utils"
	"net/http"

	"github.com/google/uuid"
)

func NewSseServer() *Server {
	return &Server{
		clients:     make(clientStore),
		topics:      make(topicStore),
		externalIds: make(externalIdMap),
	}
}

func (s *Server) ConnectClient(w http.ResponseWriter, r *http.Request, externalId string, topics []string) error {
	setSseHeaders(r)

	var flusher, ok = w.(http.Flusher)

	if !ok {
		return errors.New("streaming not supported")
	}

	var client = s.addClient(externalId, r.URL.Query().Get("id"))

	defer s.removeClient(externalId)

	err := s.subscribeToTopics(client.id, topics)

	if err != nil {
		return err
	}

	defer s.unsubscribeFromTopics(client.id, topics)

	if s.callback != nil {
		err = s.callback.OnConnection(w, r, externalId, topics, flusher)

		if err != nil {
			return err
		}
	}

	defer func() {
		if s.callback != nil {
			_ = s.callback.onDisconnection(w, r, externalId)
		}
	}()

	for messages := range client.channel {
		var formattedMessages, err = sendMultipleMessages(w, flusher, messages)

		var sent = true

		if err != nil {
			sent = false
		}

		client.messageHistory = append(client.messageHistory, clientMessageHistoryItem{
			messages:          messages,
			formattedMessages: formattedMessages,
			sent:              sent,
		})

		if s.callback != nil {
			s.callback.OnSend(err, formattedMessages, messages)
		}
	}

	return nil
}

func (s *Server) addClient(externalId string, id string) *ServerClient {
	if _, exists := s.clients[s.externalIds[externalId]]; exists {
		return utils.PtrOf(s.clients[s.externalIds[externalId]])
	}

	if _, exists := s.clients[id]; exists {
		s.externalIds[externalId] = id
		return utils.PtrOf(s.clients[s.externalIds[externalId]])
	}

	if id == "" {
		id = uuid.New().String()
	}

	client := ServerClient{
		id:         id,
		externalId: externalId,
		channel:    make(chan []Message),
	}

	s.clients[id] = client

	s.externalIds[externalId] = id

	return &client
}

func (s *Server) removeClient(externalId string) {
	var clientId = s.externalIds[externalId]

	if client, exists := s.clients[clientId]; exists {
		delete(s.externalIds, client.externalId)

		for topic, clients := range s.topics {
			newClients := make([]string, 0)

			for _, id := range clients {
				if id != clientId {
					newClients = append(newClients, id)
				}

				if len(newClients) == 0 {
					delete(s.topics, topic)
				} else {
					s.topics[topic] = newClients
				}
			}

			delete(s.clients, clientId)
		}
	}
}

func (s *Server) subscribeToTopics(clientId string, topics []string) error {
	if _, exists := s.clients[clientId]; !exists {
		return errors.New("client not found")
	}

	for _, topic := range topics {
		if _, exists := s.topics[topic]; !exists {
			s.topics[topic] = make([]string, 0)
		}

		alreadySubscribed := false

		for _, id := range s.topics[topic] {
			if id == clientId {
				alreadySubscribed = true
				break
			}
		}

		if !alreadySubscribed {
			s.topics[topic] = append(s.topics[topic], clientId)
		}
	}

	return nil
}

func (s *Server) unsubscribeFromTopics(clientId string, topics []string) {
	for _, topic := range topics {
		if clients, exists := s.topics[topic]; exists {
			newClients := make([]string, 0)

			for _, id := range clients {
				if id != clientId {
					newClients = append(newClients, id)
				}
			}

			if len(newClients) == 0 {
				delete(s.topics, topic)
			} else {
				s.topics[topic] = newClients
			}
		}
	}
}

func (s *Server) GetClientByExternalId(externalId string) (*ServerClient, error) {
	if clientId, exists := s.externalIds[externalId]; exists {
		if client, ok := s.clients[clientId]; ok {
			return &client, nil
		}
	}

	return nil, errors.New("client not found")
}

func (s *Server) SendToTopic(topic string, data MessageData) error {
	if clients, exists := s.topics[topic]; exists {
		message := Message{
			Id:    uuid.New().String(),
			Event: topic,
			Data:  data,
		}

		for _, clientId := range clients {
			if client, ok := s.clients[clientId]; ok {
				client.channel <- []Message{message}
			}
		}

		return nil
	}

	return errors.New("topic not found")
}

func (s *Server) SendToClient(clientId string, event string, data MessageData) error {
	if client, exists := s.clients[clientId]; exists {
		message := Message{
			Id:    uuid.New().String(),
			Event: event,
			Data:  data,
		}

		client.channel <- []Message{message}

		return nil
	}
	return errors.New("client not found")
}
