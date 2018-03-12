package websocket

import (
	"context"

	"go.uber.org/zap"
	ws "github.com/gorilla/websocket"
)

type WebSocketManager interface {
	Run(ctx context.Context)
	Stop()

	register(c *WebSocketClient) error
	unregister(c *WebSocketClient) error
}

type webSocketManagerImpl struct {
	clients          map[*WebSocketClient]bool
	registerChan     chan *ws.Conn
	unregisterChan   chan *WebSocketClient
	disconnectedChan chan *WebSocketClient
	finishedChan     chan bool
}

func NewWebSocketManager() *webSocketManagerImpl {
	m := &webSocketManagerImpl{
		clients:          make(map[*WebSocketClient]bool),
		registerChan:     make(chan *ws.Conn),
		unregisterChan:   make(chan *WebSocketClient),
		disconnectedChan: make(chan *WebSocketClient),
		finishedChan:     make(chan bool),
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

	count := len(m.clients)
	message, err := NewConnectionCountWebsocketMessage(count)
	if err != nil {
		logger.Errorw("Failed to build UpdateConnectionCount message")
		return err
	}

	for client := range m.clients {
		client.sendChan <- message
	}

	return nil
}

func (m *webSocketManagerImpl) unregister(c *WebSocketClient) error {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	logger.Infow("Unregister client", "uuid", c.uuid)

	if _, ok := m.clients[c]; ok {
		delete(m.clients, c)
		go func(deletedClient *WebSocketClient) {
			deletedClient.cancelFunc()
		}(c)

		count := len(m.clients)
		message, err := NewConnectionCountWebsocketMessage(count)
		if err != nil {
			logger.Errorw("Failed to build UpdateConnectionCount message", "error", err)
			return err
		}

		for client := range m.clients {
			client.sendChan <- message
		}
	}

	return nil
}

func (m *webSocketManagerImpl) run(appCtx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	logger.Info("Starting WebSocketManager")

	for {
		select {
		case client := <-m.registerChan:
			m.register(client)

		case client := <-m.unregisterChan:
			m.unregister(client)

		case <-appCtx.Done():
			for c, _ := range m.clients {
				m.unregister(c)
			}

			close(m.registerChan)
			close(m.unregisterChan)
			logger.Info("Stopped WebSocketManager")

			m.finishedChan <- true
			close(m.finishedChan)
			return
		}
	}
}

func (m *webSocketManagerImpl) Stop() <-chan bool {
	return m.finishedChan
}
