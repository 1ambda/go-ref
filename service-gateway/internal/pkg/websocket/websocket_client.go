package websocket

import (
	"time"
	"fmt"
	"encoding/json"

	"github.com/satori/go.uuid"
	ws "github.com/gorilla/websocket"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/ws_model"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/logger"
)

const (
	MessageWriteTimeout = 3 * time.Second
	MessageReadTimeout  = 5 * time.Second
	PingInterval        = 4 * time.Second
	PongTimeout         = MessageReadTimeout
)

type WebSocketMessage struct {
	Message chan [] byte
	Type    string
}

type WebSocketClient struct {
	manager    *WebSocketManager
	connection *ws.Conn
	sendChan   chan *[] byte
	closeChan  chan bool
	uuid       string
}

func NewWebSocketClient(m *WebSocketManager, conn *ws.Conn) *WebSocketClient {

	c := &WebSocketClient{
		manager:    m,
		connection: conn,
		sendChan:   make(chan *[] byte),
		closeChan:  make(chan bool),
		uuid:       uuid.NewV4().String(),
	}

	go c.run()

	return c
}

func (c *WebSocketClient) send(message *[]byte) error {
	w, err := c.connection.NextWriter(ws.TextMessage)
	if err != nil {
		return err
	}

	c.connection.SetWriteDeadline(time.Now().Add(MessageWriteTimeout))
	if _, err := w.Write(*message); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (c *WebSocketClient) close() error {
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


func buildConnectionCountMessage(count int) (*[]byte, error) {
	responseType := ws_model.WebSocketResponseHeaderResponseTypeUpdateConnectionCount
	stringified := fmt.Sprintf("%d", count)

	message := ws_model.WebSocketRealtimeResponse{
		Header: &ws_model.WebSocketResponseHeader{ResponseType: &responseType,},
		Body: &ws_model.WebSocketRealtimeResponseBody{
			Value: &stringified,
		},
	}

	serialized, err := json.Marshal(message)
	if err != nil {
		logger.Errorw("Failed to marshal UPDATE_CURRENT_CONNECTION_COUNT", "error", err)
		return nil, err
	}

	return &serialized, nil
}

func (c *WebSocketClient) validate() error {
	c.connection.SetReadDeadline(time.Now().Add(MessageReadTimeout))
	_, _, err := c.connection.ReadMessage()

	if err == nil {
		return nil
	}

	if ws.IsCloseError(err, ws.CloseGoingAway) {
		logger.Errorw("client is disconnected", "uuid", c.uuid)
		c.manager.unregisterChan <- c
		return err
	} else {
		logger.Errorw("unknown ReadMessageError", "uuid", c.uuid)
		c.manager.unregisterChan <- c
		return err
	}

	return nil
}

func (c *WebSocketClient) sendPingMessage() error {
	c.connection.SetWriteDeadline(time.Now().Add(MessageWriteTimeout))
	return c.connection.WriteMessage(ws.PingMessage, []byte{})
}

func (c *WebSocketClient) run() {
	closed := false

	pingTicker := time.NewTicker(PingInterval)

	c.connection.SetPongHandler(func(string) error {
		c.connection.SetReadDeadline(time.Now().Add(PongTimeout))
		return nil
	})

	for !closed {
		select {
		case serialized := <-c.sendChan:
			if err := c.send(serialized); err != nil {
				logger.Errorw("Failed to send message to client", "uuid", c.uuid)
				c.manager.unregisterChan <- c
			}

		case <-pingTicker.C:
			logger.Infow("WebSocketClient Ping", "uuid", c.uuid)
			if err := c.sendPingMessage(); err != nil {
				logger.Errorw("Failed to ping message to client", "uuid", c.uuid)
				c.manager.unregisterChan <- c
			}

		case <-c.closeChan:
			closed = true
			break
		}
	}

	c.close()
	close(c.sendChan)
	close(c.closeChan)
}
