package websocket

import (
	"go.uber.org/zap"
	"sync"
)

type WebSocketManager struct {
	clients          map[*WebSocketClient]bool
	registerChan     chan *WebSocketClient
	unregisterChan   chan *WebSocketClient
	disconnectedChan chan *WebSocketClient
	closeChan        chan bool
	finishedChan     chan bool
	lock             sync.Mutex
}

func NewWebSocketManager() *WebSocketManager {
	m := &WebSocketManager{
		clients:          make(map[*WebSocketClient]bool),
		registerChan:     make(chan *WebSocketClient),
		unregisterChan:   make(chan *WebSocketClient),
		disconnectedChan: make(chan *WebSocketClient),
		closeChan:        make(chan bool),
		finishedChan:     make(chan bool),
	}

	go m.run()

	return m
}

func (m *WebSocketManager) register(c *WebSocketClient) error {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	logger.Infow("Register client", "uuid", c.uuid)

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
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	logger.Infow("Unregister client", "uuid", c.uuid)

	if _, ok := m.clients[c]; ok {
		delete(m.clients, c)
		go func(deletedClient *WebSocketClient) {
			deletedClient.closeChan <- true
		}(c)

		count := len(m.clients)
		message, err := NewConnectionCountWebsocketMessage(count)
		if err != nil {
			logger.Errorw("Failed to build UpdateConnectionCount message", "error", err)
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
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	logger.Info("Starting WebSocketManager")

	done := false
	for !done {
		select {
		case client := <-m.registerChan:
			m.register(client)

		case client := <-m.unregisterChan:
			m.unregister(client)

		case <-m.closeChan:
			for c, _ := range m.clients {
				m.unregister(c)
			}

			done = true
			break
		}
	}

	close(m.registerChan)
	close(m.unregisterChan)
	logger.Info("Stopped WebSocketManager")
	m.finishedChan <- true
}

func (m *WebSocketManager) Stop() {
	m.closeChan <- true
	<-m.finishedChan
	close(m.finishedChan)
}
