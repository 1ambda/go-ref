package websocket

import (
	"context"

	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
	"fmt"
)

type WebSocketManager interface {
	Broadcast(message *WebSocketMessage)
	SubscribeConnectionCount() <-chan string // subscribed by DistributedClient
	Stop() <-chan bool
}

type webSocketManagerImpl struct {
	clients               map[*WebSocketClient]bool
	broadcastMessageCache map[string]*WebSocketMessage // used to cache recent values

	registerChan     chan *ws.Conn
	unregisterChan   chan *WebSocketClient
	disconnectedChan chan *WebSocketClient
	broadcastChan    chan *WebSocketMessage
	finishedChan     chan bool

	wsConnCountChan chan string
}

func NewWebSocketManager() *webSocketManagerImpl {
	m := &webSocketManagerImpl{
		clients:               make(map[*WebSocketClient]bool),
		broadcastMessageCache: make(map[string]*WebSocketMessage),

		registerChan:          make(chan *ws.Conn),
		unregisterChan:        make(chan *WebSocketClient),
		disconnectedChan:      make(chan *WebSocketClient),
		broadcastChan:         make(chan *WebSocketMessage),
		finishedChan:          make(chan bool),

		wsConnCountChan: make(chan string),
	}

	return m
}

func (m *webSocketManagerImpl) register(conn *ws.Conn) error {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

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

func (m *webSocketManagerImpl) unregister(c *WebSocketClient) error {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

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

func (m *webSocketManagerImpl) run(appCtx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	logger.Info("Starting WebSocketManager")

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

			logger.Info("Stopped WebSocketManager")
			m.finishedChan <- true
			close(m.finishedChan)
			return
		}
	}
}

func (m *webSocketManagerImpl) SubscribeConnectionCount() <-chan string {
	return m.wsConnCountChan
}

func (m *webSocketManagerImpl) Broadcast(message *WebSocketMessage) {
	m.broadcastChan <- message
}

func (m *webSocketManagerImpl) Stop() <-chan bool {
	return m.finishedChan
}
