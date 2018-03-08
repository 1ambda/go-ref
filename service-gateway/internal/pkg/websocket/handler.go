package websocket

import (
	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Configure(mux *http.ServeMux, m *WebSocketManager) {
	mux.HandleFunc("/endpoint", func(res http.ResponseWriter, req *http.Request) {
		log, _ := zap.NewProduction()
		defer log.Sync() // flushes buffer, if any
		logger := log.Sugar()

		conn, err := upgrader.Upgrade(res, req, nil)
		if err != nil {
			logger.Errorw("Failed to get WS connection", "error", err)
			return
		}

		// register a client
		client := NewWebSocketClient(m, conn)
		m.registerChan <- client
	})
}
