package websocket

import (
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/logger"
	"sync"
)

type WebSocketManager struct {
	clients          map[*WebSocketClient]bool
	registerChan     chan *WebSocketClient
	unregisterChan   chan *WebSocketClient
	disconnectedChan chan *WebSocketClient
	broadcastChan    chan [] byte
	lock             sync.Mutex
}

func NewWebSocketManager() *WebSocketManager {
	m := &WebSocketManager{
		clients:          make(map[*WebSocketClient]bool),
		registerChan:     make(chan *WebSocketClient),
		unregisterChan:   make(chan *WebSocketClient),
		disconnectedChan: make(chan *WebSocketClient),
		broadcastChan:    make(chan []byte),
	}

	go m.run()

	return m
}

func (m *WebSocketManager) register(c *WebSocketClient) error {
	logger.Infow("Requesting client registration", "uuid", c.uuid)

	m.clients[c] = true

	count := len(m.clients)
	message, err := NewConnectionCountWebsocketMessage(count)
	if err != nil {
		logger.Errorw("Failed to build UpdateConnectionCount message")
		return err
	}

	for client := range m.clients {
		m.signalToSendMessage(client, message)
	}

	return nil
}

func (m *WebSocketManager) unregister(c *WebSocketClient) error {
	logger.Infow("Requesting client removal", "uuid", c.uuid)

	if _, ok := m.clients[c]; ok {
		delete(m.clients, c)
		go func(deletedClient *WebSocketClient) {
			deletedClient.closeChan <- true
		}(c)

		count := len(m.clients)
		message, err := NewConnectionCountWebsocketMessage(count)
		if err != nil {
			logger.Errorw("Failed to build UpdateConnectionCount message")
			return err
		}

		for client := range m.clients {
			m.signalToSendMessage(client, message)
		}
	}

	return nil
}

func (m *WebSocketManager) signalToSendMessage(c *WebSocketClient, msg *WebSocketMessage) {
	c.sendChan <- msg
}

func (m *WebSocketManager) run() {
	// TODO: close

	for {
		select {
		case client := <-m.registerChan:
			m.register(client)

		case client := <-m.unregisterChan:
			// TODO: broadcast
			m.unregister(client)
		}
	}

	close(m.registerChan)
	close(m.unregisterChan)
}
