package websocket

import (
	"context"
	"fmt"
	"time"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
	"github.com/1ambda/go-ref/service-gateway/internal/model"
	"github.com/jinzhu/gorm"
)

type Manager interface {
	Broadcast(message *Message)
	SubscribeConnectionCount() <-chan string // subscribed by DistributedClient
	Stop() <-chan bool
}

type managerImpl struct {
	clients               map[*Client]bool
	broadcastMessageCache map[string]*Message // used to cache recent values

	registerChan     chan *Client
	unregisterChan   chan *Client
	disconnectedChan chan *Client
	broadcastChan    chan *Message
	finishedChan     chan bool

	wsConnCountChan chan string

	db *gorm.DB
}

func NewManager(db *gorm.DB) *managerImpl {
	m := &managerImpl{
		clients:               make(map[*Client]bool),
		broadcastMessageCache: make(map[string]*Message),

		registerChan:     make(chan *Client),
		unregisterChan:   make(chan *Client),
		disconnectedChan: make(chan *Client),
		broadcastChan:    make(chan *Message),
		finishedChan:     make(chan bool),

		wsConnCountChan: make(chan string),
		db:              db,
	}

	return m
}

func (m *managerImpl) register(c *Client) error {
	logger := config.GetLogger()

	ctx, cancel := context.WithCancel(context.Background())
	c.cancelFunc = cancel

	sessionID := c.getSessionID()
	websocketID := c.getWebsocketID()

	logger.Infow("Register client", "session_id", sessionID, "websocket_id", websocketID)

	m.clients[c] = true
	go c.run(ctx)

	// send cached broadcast messages
	for _, message := range m.broadcastMessageCache {
		c.sendChan <- message
	}

	count := fmt.Sprintf("%d", len(m.clients))
	m.wsConnCountChan <- count

	return nil
}

func (m *managerImpl) unregister(c *Client) {
	logger := config.GetLogger()

	if _, ok := m.clients[c]; !ok {
		return // ignore a request for an invalid client
	}

	sessionID := c.getSessionID()
	websocketID := c.getWebsocketID()

	logger.Infow("Unregister client", "session_id", sessionID, "websocket_id", websocketID)

	// delete from the in-memory client map
	delete(m.clients, c)

	// send message to etcd client
	count := fmt.Sprintf("%d", len(m.clients))
	m.wsConnCountChan <- count

	// update websocket history table for this client
	m.updateWebsocketHistory(c, false)

	// call client shutdown hook
	go func(deletedClient *Client) {
		deletedClient.cancelFunc()
	}(c)

	return
}

func (m *managerImpl) updateWebsocketHistory(c *Client, isShutdown bool) {
	logger := config.GetLogger()
	sessionID := c.getSessionID()
	websocketID := c.getWebsocketID()

	// get closeReason
	closeReason := c.getCloseReason()
	if closeReason == "" {
		closeReason = config.WsCloseReasonUnknown
	}
	if isShutdown {
		closeReason = config.WsCloseReasonServerShutdown
	}

	// update websocket_history columns
	record := model.WebsocketHistory{}
	result := m.db.Model(&record).
		Where("websocket_id = ?", websocketID).
		Updates(map[string]interface{}{
		"close_reason": closeReason,
		"closed_at":    time.Now().UTC(),
	})

	if result.Error != nil {
		logger.Errorw("Failed to update WebsocketHistory record due to unknown error",
			"session_id", sessionID, "websocket_id", websocketID, "error", result.Error)
		return
	}

	if result.RowsAffected < 1 {
		logger.Errorw("Failed to find WebsocketHistory record before updating",
			"session_id", sessionID, "websocket_id", websocketID, "error", result.Error)
		return
	}
}

func (m *managerImpl) run(appCtx context.Context) {
	logger := config.GetLogger()
	logger.Info("Starting Manager")

	for {
		select {
		case message := <-m.broadcastChan:
			// update cache. this will be used for newly joined clients
			m.broadcastMessageCache[message.event] = message

			for client := range m.clients {
				client.sendChan <- message
			}

		case registration := <-m.registerChan:
			m.register(registration)

		case client := <-m.unregisterChan:
			m.unregister(client)

		case <-appCtx.Done():
			for client := range m.clients {
				m.updateWebsocketHistory(client, true)
			}

			close(m.registerChan)
			close(m.unregisterChan)
			close(m.broadcastChan)

			logger.Info("Stopped Manager")
			m.finishedChan <- true
			close(m.finishedChan)
			return
		}
	}
}

func (m *managerImpl) SubscribeConnectionCount() <-chan string {
	return m.wsConnCountChan
}

func (m *managerImpl) Broadcast(message *Message) {
	m.broadcastChan <- message
}

func (m *managerImpl) Stop() <-chan bool {
	return m.finishedChan
}
