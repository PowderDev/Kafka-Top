package websockets

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Hub struct {
	clients          map[*WsClient]bool
	OutgoingMessages chan Message
	register         chan *WsClient
	unregister       chan *WsClient
}

var clientsMutex = &sync.Mutex{}

func NewHub() *Hub {
	return &Hub{
		OutgoingMessages: make(chan Message),
		register:         make(chan *WsClient),
		unregister:       make(chan *WsClient),
		clients:          make(map[*WsClient]bool),
	}
}

func (h *Hub) broadcastData(message Message) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}

	clientsMutex.Lock()
	for client := range h.clients {
		client.writeMessage(messageBytes)
	}
	clientsMutex.Unlock()
}

func (h *Hub) HandleKafkaMessage(message *kafka.Message) {
	var parsedValue interface{}
	err := json.Unmarshal(message.Value, &parsedValue)
	if err != nil {
		log.Println(err)
		return
	}

	h.OutgoingMessages <- Message{
		Key: string(message.Key),
		Data: map[string]interface{}{
			"value": parsedValue,
			"ts":    message.Timestamp.UnixMilli(),
		},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			clientsMutex.Lock()
			h.clients[client] = true
			clientsMutex.Unlock()

		case client := <-h.unregister:
			clientsMutex.Lock()
			delete(h.clients, client)
			clientsMutex.Unlock()

		case message := <-h.OutgoingMessages:
			h.broadcastData(message)
		}
	}
}

func (h *Hub) Stop() {
	for client := range h.clients {
		client.conn.Close()
	}
}
