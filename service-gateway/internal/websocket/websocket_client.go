package websocket

import (
	"context"
	"sync"
	"time"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
	ws "github.com/gorilla/websocket"
)

const (
	MessageWriteTimeout = 300 * time.Millisecond
	MessageReadTimeout  = 3 * time.Second // should be greater than `PingInterval`
	PingInterval        = 2 * time.Second
	PongTimeout         = MessageReadTimeout
	MessagePopInterval  = 1 * time.Second // delay until pop the next message from buffer
)

type Client struct {
	manager     *managerImpl
	connection  *ws.Conn
	sendChan    chan *Message
	buffer      []*Message
	websocketID string
	sessionID   string
	cancelFunc  context.CancelFunc
	closeReason string

	lock sync.RWMutex
}

func (c *Client) getWebsocketID() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.websocketID
}

func (c *Client) getSessionID() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.sessionID
}

func (c *Client) getCloseReason() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.closeReason
}

func (c *Client) setCloseReason(reason string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.closeReason = reason
}

func NewClient(m *managerImpl, conn *ws.Conn, sessionID string, websocketID string) *Client {

	c := &Client{
		manager:     m,
		connection:  conn,
		sendChan:    make(chan *Message),
		buffer:      make([]*Message, 0),
		sessionID:   sessionID,
		websocketID: websocketID,
		// cancelFunc is configured in the websocket manager
	}

	return c
}

func (c *Client) send(message *Message) error {
	logger := config.GetLogger()
	websocketID := c.getWebsocketID()

	w, err := c.connection.NextWriter(ws.TextMessage)
	if err != nil {
		logger.Errorw("Failed to get next writer", "websocket_id", websocketID, "error", err)
		return err
	}
	defer w.Close()

	logger.Debugw("Sending websocket message to client",
		"websocket_id", websocketID, "event", message.event)

	c.connection.SetWriteDeadline(time.Now().Add(MessageWriteTimeout))
	if _, err := w.Write(*message.content); err != nil {
		return err
	}

	return nil
}

func (c *Client) close() error {
	logger := config.GetLogger()
	websocketID := c.getWebsocketID()

	logger.Infow("Closing client", "websocket_id", websocketID)

	if err := c.connection.WriteMessage(ws.CloseMessage, []byte{}); err != nil {
		logger.Warnw("Failed to send `CloseMessage`", "websocket_id", websocketID)
	}

	if err := c.connection.Close(); err != nil {
		logger.Warnw("Failed to close client", "websocketID", websocketID)
		return err
	}

	return nil
}

func (c *Client) sendPingMessage() error {
	c.connection.SetWriteDeadline(time.Now().Add(MessageWriteTimeout))
	return c.connection.WriteMessage(ws.PingMessage, []byte{})
}

func (c *Client) SendErrorMessage(err error, errorType string, code int64) {
	logger := config.GetLogger()
	message, serializeErr := NewErrorMessage(err, errorType, code)
	if serializeErr != nil {
		logger.Errorw("Failed to create websocket error message", "error", serializeErr)
		return
	}

	sendErr := c.send(message)
	if sendErr != nil {
		logger.Errorw("Failed to send websocket error message", "error", sendErr)
	}
}

func (c *Client) run(ctx context.Context) {
	logger := config.GetLogger()

	pingTicker := time.NewTicker(PingInterval)
	messagePopTicker := time.NewTicker(MessagePopInterval)

	c.connection.SetPongHandler(func(string) error {
		c.connection.SetReadDeadline(time.Now().Add(PongTimeout))
		return nil
	})

	for {
		select {
		case message := <-c.sendChan:
			c.buffer = append(c.buffer, message)

		case <-pingTicker.C:
			if err := c.sendPingMessage(); err != nil {
				websocketID := c.getWebsocketID()
				logger.Warnw("Failed to ping message to client",
					"websocket_id", websocketID)
				c.closeReason = config.WsCloseFailureClientDisconnected
				c.manager.unregisterChan <- c
			}

		case <-messagePopTicker.C:
			if len(c.buffer) == 0 {
				continue
			}

			for _, message := range c.buffer {
				if err := c.send(message); err != nil {
					// don't write log, client is disconnected
					websocketID := c.getWebsocketID()
					logger.Warnw("Failed to send message to client. Closing this client",
						"websocket_id", websocketID)
					c.setCloseReason(config.WsCloseReasonMessageSendFailure)
					c.manager.unregisterChan <- c
					break
				}
			}

			c.buffer = nil

		case <-ctx.Done():
			c.close()
			close(c.sendChan)
			return
		}
	}
}
