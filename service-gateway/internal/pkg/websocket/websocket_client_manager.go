package websocket

import (
	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
	"encoding/json"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/ws_model"
	"fmt"
	"time"
)

type WebSocketClientManager struct {
	clients    map[*WebSocketClient]bool
	register   chan *WebSocketClient
	unregister chan *WebSocketClient
	broadcast  chan [] byte
}

type WebSocketClient struct {
	manager    *WebSocketClientManager
	connection *ws.Conn
}

func NewWebSocketClientManager() *WebSocketClientManager {
	m := &WebSocketClientManager{
		clients:    make(map[*WebSocketClient]bool),
		register:   make(chan *WebSocketClient),
		unregister: make(chan *WebSocketClient),
		broadcast:  make(chan []byte),
	}

	go func() {
		m.run()
	}()

	return m
}

func NewWebSocketClient(m *WebSocketClientManager, conn *ws.Conn) *WebSocketClient {
	return &WebSocketClient{
		manager: m, connection: conn,
	}
}

func (c *WebSocketClient) register() error {
	if c == nil {
		return nil
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Info("Registering client")

	c.manager.register <- c

	return nil
}

func (c *WebSocketClient) unregister() error {
	if c == nil {
		return nil
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Info("Registering client")

	c.manager.unregister <- c

	return nil
}

func (c *WebSocketClient) send(message *[]byte) error {
	if c == nil {
		return nil
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Info("Sending a message to client")

	w, err := c.connection.NextWriter(ws.TextMessage)
	if err != nil {
		sugar.Errorw("Failed to get writer", "error", err)
		return nil
	}

	w.Write(*message)

	if err := w.Close(); err != nil {
		sugar.Errorw("Failed to close writer", "error", err)
		return nil
	}

	return nil
}

func (c *WebSocketClient) close() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	if err := c.connection.WriteMessage(ws.CloseMessage, []byte{}); err != nil {
		sugar.Errorw("Failed to send `CloseMessage`", "error", err)
		return
	}

	if err := c.connection.Close(); err != nil {
		sugar.Errorw("Failed to close a websocket client", "error", err)
		return
	}

	return
}

func (c *WebSocketClient) sendUpdateConnectionCountMessage(clientsCount int) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

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
		sugar.Errorw("Failed to marshal UPDATE_CURRENT_CONNECTION_COUNT", "error", err)
		return err
	}

	c.send(&serialized)

	return nil
}

func (m *WebSocketClientManager) run() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	for {
		select {
		case client := <-m.register:
			sugar.Info("Registering a client in manager")
			m.clients[client] = true

			for client := range m.clients {
				client.sendUpdateConnectionCountMessage(len(m.clients))
			}

		case client := <-m.unregister:
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				client.close()
			}

		case message := <-m.broadcast:
			for client := range m.clients {
				client.send(&message)
			}

		case <- time.After(3 * time.Second):
			for client := range m.clients {
				_, _, err := client.connection.ReadMessage()
				if err != nil {
					sugar.Errorw("ReadMessage", "error", err)
					delete(m.clients, client)
					client.close()
				}
			}
		}
	}
}
