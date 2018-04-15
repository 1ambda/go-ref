package websocket

import (
	"context"
	"fmt"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
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

func (m *managerImpl) register(client *Client) error {
	logger := config.GetLogger()

	ctx, cancel := context.WithCancel(context.Background())
	client.cancelFunc = cancel
	logger.Infow("Register client", "websocketID", client.websocketID)

	m.clients[client] = true
	go client.run(ctx)

	// send cached broadcast messages
	for _, message := range m.broadcastMessageCache {
		client.sendChan <- message
	}

	count := fmt.Sprintf("%d", len(m.clients))
	m.wsConnCountChan <- count

	return nil
}

func (m *managerImpl) unregister(c *Client) error {
	logger := config.GetLogger()
	logger.Infow("Unregister client", "websocketID", c.websocketID)

	if _, ok := m.clients[c]; !ok {
		return nil // ignore a request for an invalid client
	}

	delete(m.clients, c)
	go func(deletedClient *Client) {
		deletedClient.cancelFunc()
	}(c)

	count := fmt.Sprintf("%d", len(m.clients))
	m.wsConnCountChan <- count

	return nil
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
