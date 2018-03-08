package websocket

import (
	"net/http"
	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
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
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		sugar := logger.Sugar()

		conn, err := upgrader.Upgrade(res, req, nil)
		if err != nil {
			sugar.Errorw("Failed to get WS connection", "error", err)
			return
		}

		// register a client
		client := NewWebSocketClient(m, conn)
		m.registerChan <- client
	})
}
