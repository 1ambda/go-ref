package websocket

import (
	ws "github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

const (
	MessageWriteTimeout = 300 * time.Millisecond
	MessageReadTimeout  = 3 * time.Second // should be greater than `PingInterval`
	PingInterval        = 2 * time.Second
	PongTimeout         = MessageReadTimeout
	MessagePopInterval  = 1 * time.Second // delay until pop the next message from buffer
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
	defer log.Sync()
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
	defer log.Sync()
	logger := log.Sugar()

	logger.Infow("Closing client", "uuid", c.uuid)

	if err := c.connection.WriteMessage(ws.CloseMessage, []byte{}); err != nil {
		logger.Warnw("Failed to send `CloseMessage`", "uuid", c.uuid)
	}

	if err := c.connection.Close(); err != nil {
		logger.Warnw("Failed to close client", "uuid", c.uuid)
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
	defer log.Sync()
	logger := log.Sugar()

	done := false
	pingTicker := time.NewTicker(PingInterval)
	messagePopTicker := time.NewTicker(MessagePopInterval)

	c.connection.SetPongHandler(func(string) error {
		c.connection.SetReadDeadline(time.Now().Add(PongTimeout))
		return nil
	})

	for !done {
		select {
		case message := <-c.sendChan:
			c.buffer = append(c.buffer, message)

		case <-pingTicker.C:
			if err := c.sendPingMessage(); err != nil {
				logger.Warnw("Failed to ping message to client", "uuid", c.uuid)
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
				// don't write log, client is disconnected
				logger.Warnw("Failed to send message to client", "uuid", c.uuid)
				c.manager.unregisterChan <- c
			}
			c.isSending = false

		case <-c.closeChan:
			done = true
			break
		}
	}

	c.close()
	close(c.sendChan)
	close(c.closeChan)
}
