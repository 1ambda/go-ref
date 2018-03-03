package websocketservice

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"github.com/1ambda/go-ref/service-gateway/pkg/grpc"
	"fmt"
	//"encoding/json"
	"github.com/golang/protobuf/proto"
)

type WebSocketClientManager struct {
	clients    map[*WebSocketClient]bool
	register   chan *WebSocketClient
	unregister chan *WebSocketClient
	broadcast  chan [] byte
}

type WebSocketClient struct {
	manager    *WebSocketClientManager
	connection *websocket.Conn
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

func NewWebSocketClient(m *WebSocketClientManager, conn *websocket.Conn) *WebSocketClient {
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

	w, err := c.connection.NextWriter(websocket.TextMessage)
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
	c.connection.WriteMessage(websocket.CloseMessage, []byte{})
	c.connection.Close()
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

			message := &grpc.WebSocketRealtimeResponse{
				Header: &grpc.WebSocketResponseHeader{
					EventType: grpc.WebSocketResponseHeader_UPDATE_CURRENT_CONNECTION_COUNT,
				},
				Body: &grpc.WebSocketRealtimeResponseBody{
					Value: fmt.Sprintf("%d", len(m.clients)),
				},
			}

			//serialized, err := json.Marshal(message)
			serialized, err := proto.Marshal(message)
			if err != nil {
				sugar.Errorw("Failed to marshal UPDATE_CURRENT_CONNECTION_COUNT", "error", err)
				continue
			}

			sugar.Info("Broadcasting a message to all clients")
			for client := range m.clients {
				client.send(&serialized)
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
		}
	}
}
