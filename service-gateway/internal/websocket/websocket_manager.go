package websocket

import (
	"context"
	"fmt"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
	ws "github.com/gorilla/websocket"
)

type Manager interface {
	Broadcast(message *Message)
	SubscribeConnectionCount() <-chan string // subscribed by DistributedClient
	Stop() <-chan bool
}

type managerImpl struct {
	clients               map[*WebSocketClient]bool
	broadcastMessageCache map[string]*Message // used to cache recent values

	registerChan     chan *ws.Conn
	unregisterChan   chan *WebSocketClient
	disconnectedChan chan *WebSocketClient
	broadcastChan    chan *Message
	finishedChan     chan bool

	wsConnCountChan chan string
}

func NewManager() *managerImpl {
	m := &managerImpl{
		clients:               make(map[*WebSocketClient]bool),
		broadcastMessageCache: make(map[string]*Message),

		registerChan:     make(chan *ws.Conn),
		unregisterChan:   make(chan *WebSocketClient),
		disconnectedChan: make(chan *WebSocketClient),
		broadcastChan:    make(chan *Message),
		finishedChan:     make(chan bool),

		wsConnCountChan: make(chan string),
	}

	return m
}

func (m *managerImpl) register(conn *ws.Conn) error {
	logger := config.GetLogger()

	ctx, cancel := context.WithCancel(context.Background())
	client := NewWebSocketClient(m, conn, cancel)
	logger.Infow("Register client", "uuid", client.uuid)

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

func (m *managerImpl) unregister(c *WebSocketClient) error {
	logger := config.GetLogger()
	logger.Infow("Unregister client", "uuid", c.uuid)

	if _, ok := m.clients[c]; !ok {
		return nil // ignore a request for an invalid client
	}

	delete(m.clients, c)
	go func(deletedClient *WebSocketClient) {
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

		case client := <-m.registerChan:
			m.register(client)

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
