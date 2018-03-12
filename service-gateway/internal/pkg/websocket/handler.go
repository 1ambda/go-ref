package websocket

import (
	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"context"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Configure(appCtx context.Context, mux *http.ServeMux) *webSocketManagerImpl {
	// start websocket manager
	webSocketManager := NewWebSocketManager()
	go webSocketManager.run(appCtx)

	// setup endpoint
	mux.HandleFunc("/endpoint", func(res http.ResponseWriter, req *http.Request) {
		log, _ := zap.NewProduction()
		defer log.Sync()
		logger := log.Sugar()

		conn, err := upgrader.Upgrade(res, req, nil)
		if err != nil {
			logger.Errorw("Failed to get WS connection", "error", err)
			return
		}

		// register a client
		webSocketManager.registerChan <- conn
	})

	return webSocketManager
}
