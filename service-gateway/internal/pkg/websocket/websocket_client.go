package websocket

import (
	ws "github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

const (
	MessageWriteTimeout = 2 * time.Second
	MessageReadTimeout  = 5 * time.Second
	MessagePopInterval  = 1 * time.Second
	PingInterval        = 4 * time.Second
	PongTimeout         = MessageReadTimeout
)

type WebSocketClient struct {
	manager    *WebSocketManager
	connection *ws.Conn
	sendChan   chan *WebSocketMessage
	closeChan  chan bool
	buffer     []*WebSocketMessage
	uuid       string
	isSending  bool
}

func NewWebSocketClient(m *WebSocketManager, conn *ws.Conn) *WebSocketClient {

	c := &WebSocketClient{
		manager:    m,
		connection: conn,
		sendChan:   make(chan *WebSocketMessage),
		closeChan:  make(chan bool),
		buffer:     make([]*WebSocketMessage, 0),
		uuid:       uuid.NewV4().String(),
		isSending:  false,
	}

	go c.run()

	return c
}

func (c *WebSocketClient) send(message *WebSocketMessage) error {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

	w, err := c.connection.NextWriter(ws.TextMessage)
	if err != nil {
		logger.Errorw("Failed to get next writer", "uuid", c.uuid, "error", err)
		return err
	}
	defer w.Close()

	logger.Debugw("Sending websocket message to client",
		"uuid", c.uuid, "event", message.event)

	c.connection.SetWriteDeadline(time.Now().Add(MessageWriteTimeout))
	if _, err := w.Write(*message.content); err != nil {
		return err
	}

	return nil
}

func (c *WebSocketClient) close() error {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

	logger.Infow("Closing client", "uuid", c.uuid)

	if err := c.connection.WriteMessage(ws.CloseMessage, []byte{}); err != nil {
		logger.Errorw("Failed to send `CloseMessage`", "uuid", c.uuid, "error", err)
	}

	if err := c.connection.Close(); err != nil {
		logger.Errorw("Failed to close client", "uuid", c.uuid, "error", err)
		return err
	}

	return nil
}

func (c *WebSocketClient) sendPingMessage() error {
	c.connection.SetWriteDeadline(time.Now().Add(MessageWriteTimeout))
	return c.connection.WriteMessage(ws.PingMessage, []byte{})
}

func (c *WebSocketClient) run() {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

	closed := false
	pingTicker := time.NewTicker(PingInterval)
	messagePopTicker := time.NewTicker(MessagePopInterval)

	c.connection.SetPongHandler(func(string) error {
		c.connection.SetReadDeadline(time.Now().Add(PongTimeout))
		return nil
	})

	for !closed {
		select {
		case message := <-c.sendChan:
			c.buffer = append(c.buffer, message)

		case <-pingTicker.C:
			if err := c.sendPingMessage(); err != nil {
				logger.Errorw("Failed to ping message to client", "uuid", c.uuid)
				c.manager.unregisterChan <- c
			}

		case <-messagePopTicker.C:
			if len(c.buffer) == 0 || c.isSending {
				continue
			}

			message := c.buffer[0]
			c.buffer = c.buffer[1:]

			c.isSending = true
			if err := c.send(message); err != nil {
				logger.Errorw("Failed to send message to client", "uuid", c.uuid)
				c.manager.unregisterChan <- c
			}
			c.isSending = false

		case <-c.closeChan:
			closed = true
			break
		}
	}

	c.close()
	close(c.sendChan)
	close(c.closeChan)
}
