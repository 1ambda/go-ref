package websocket

import (
	ws "github.com/gorilla/websocket"
	"encoding/json"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/ws_model"
	"fmt"
	"time"
	"sync"
	"github.com/satori/go.uuid"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/logger"
)

type WebSocketManager struct {
	clients        map[*WebSocketClient]bool
	registerChan   chan *WebSocketClient
	deregisterChan chan *WebSocketClient
	broadcastChan  chan [] byte
}

const (
	MessageReadTimeout  = 3 * time.Second
	MessageWriteTimeout = 3 * time.Second
	ClientCheckInterval = 2 * time.Second
)

type WebSocketClient struct {
	manager    *WebSocketManager
	connection *ws.Conn
	uuid       string
	lock       sync.Mutex
}

func NewWebSocketManager() *WebSocketManager {
	m := &WebSocketManager{
		clients:        make(map[*WebSocketClient]bool),
		registerChan:   make(chan *WebSocketClient),
		deregisterChan: make(chan *WebSocketClient),
		broadcastChan:  make(chan []byte),
	}

	go func() {
		m.run()
	}()

	return m
}

func NewWebSocketClient(m *WebSocketManager, conn *ws.Conn) *WebSocketClient {
	return &WebSocketClient{
		manager: m, connection: conn, uuid: uuid.NewV4().String(),
	}
}

func (m *WebSocketManager) register(c *WebSocketClient) error {
	if c == nil {
		return nil
	}

	logger.Infow("Requesting client registration", "uuid", c.uuid)

	m.registerChan <- c

	return nil
}

func (m *WebSocketManager) unregister(c *WebSocketClient) error {
	if m == nil {
		return nil
	}

	logger.Infow("Requesting client deregisteration", "uuid", c.uuid)

	m.deregisterChan <- c

	return nil
}

func (c *WebSocketClient) send(message *[]byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	logger.Infow("Sending message to client", "uuid", c.uuid)

	w, err := c.connection.NextWriter(ws.TextMessage)
	if err != nil {
		logger.Errorw("Failed to get writer", "error", err)
		return err
	}

	c.connection.SetWriteDeadline(time.Now().Add(MessageWriteTimeout))
	if _, err := w.Write(*message); err != nil {
		logger.Errorw("Failed to write", "error", err)
		return err
	}

	if err := w.Close(); err != nil {
		logger.Errorw("Failed to close writer", "error", err)
		return err
	}

	return nil
}

func (c *WebSocketClient) close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	logger.Infow("Closing client", "uuid", c.uuid)

	if err := c.connection.WriteMessage(ws.CloseMessage, []byte{}); err != nil {
		logger.Errorw("Failed to send `CloseMessage`", "uuid", c.uuid, "error", err)
		return err
	}

	if err := c.connection.Close(); err != nil {
		logger.Errorw("Failed to close client", "uuid", c.uuid, "error", err)
		return err
	}

	return nil
}

func (c *WebSocketClient) sendUpdateConnectionCountMessage(clientsCount int) error {
	responseType := ws_model.WebSocketResponseHeaderResponseTypeUpdateConnectionCount
	count := fmt.Sprintf("%d", clientsCount)

	message := ws_model.WebSocketRealtimeResponse{
		Header: &ws_model.WebSocketResponseHeader{ResponseType: &responseType,},
		Body: &ws_model.WebSocketRealtimeResponseBody{
			Value: &count,
		},
	}

	serialized, err := json.Marshal(message)
	if err != nil {
		logger.Errorw("Failed to marshal UPDATE_CURRENT_CONNECTION_COUNT", "error", err)
		return err
	}

	c.send(&serialized)

	return nil
}

func (c *WebSocketClient) validate() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.connection.SetReadDeadline(time.Now().Add(MessageReadTimeout))
	_, _, err := c.connection.ReadMessage()

	if err == nil {
		return nil
	}

	if ws.IsCloseError(err, ws.CloseGoingAway) {
		logger.Errorw("client is disconnected", "uuid", c.uuid)
		c.manager.deregisterChan <- c
		return err
	} else {
		logger.Errorw("unknown ReadMessageError", "uuid", c.uuid)
		c.manager.deregisterChan <- c
		return err
	}

	return nil
}

func (m *WebSocketManager) run() {
	for {
		select {
		case client := <-m.registerChan:
			logger.Info("Registering a client in manager")

			m.clients[client] = true

			for client := range m.clients {
				client.sendUpdateConnectionCountMessage(len(m.clients))
			}

		case client := <-m.deregisterChan:
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				client.close()
			}

			for client := range m.clients {
				client.sendUpdateConnectionCountMessage(len(m.clients))
			}

		case message := <-m.broadcastChan:
			for client := range m.clients {
				client.send(&message)
			}

		case <-time.After(ClientCheckInterval):
			logger.Infow("websocket client liveness check", "count", len(m.clients))

			for client := range m.clients {
				client.validate()
			}
		}
	}
}
