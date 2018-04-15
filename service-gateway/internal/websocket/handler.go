package websocket

import (
	"context"
	"net/http"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Configure(appCtx context.Context, mux *http.ServeMux) *managerImpl {
	// start websocket manager
	webSocketManager := NewManager()
	go webSocketManager.run(appCtx)

	// setup endpoint
	mux.HandleFunc("/endpoint", func(res http.ResponseWriter, req *http.Request) {
		logger := config.GetLogger()

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
