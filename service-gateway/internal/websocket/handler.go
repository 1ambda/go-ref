package websocket

import (
	"context"
	"errors"
	"net/http"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
	"github.com/1ambda/go-ref/service-gateway/internal/model"
	ws "github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientRegistrationRequest struct {
	Conn        *ws.Conn // connection
	SessionID   string   // session id
	WebsocketID string   // websocket websocketID
}

func Configure(appCtx context.Context, mux *http.ServeMux, db *gorm.DB) *managerImpl {
	// start websocket manager
	manager := NewManager(db)
	go manager.run(appCtx)

	// setup endpoint
	mux.HandleFunc("/endpoint", func(res http.ResponseWriter, req *http.Request) {
		logger := config.GetLogger()

		conn, err := upgrader.Upgrade(res, req, nil)
		if err != nil {
			logger.Errorw("Failed to get WS connection", "error", err)
			return
		}

		// get session id
		sessionID, err := getSessionCookieForWs(req)
		if err != nil {
			logger.Errorw("Failed to get Session Cookie", "error", err)
			return
		}

		websocketID := uuid.NewV4().String()

		client := NewClient(manager, conn, sessionID, websocketID)

		// create websocket history record
		record := model.WebsocketHistory{}
		record.NewWebSocketHistory(sessionID, websocketID)
		if err := db.Create(&record).Error; err != nil {
			message, err1 := NewErrorMessage(err, 500)
			if err1 != nil {
				logger.Errorw("Failed to create websocket error message", "error", err1)
				return
			}

			client.send(message)
			client.close()
			return
		}

		// register a client
		manager.registerChan <- client
	})

	return manager
}

func getSessionCookieForWs(req *http.Request) (string, error) {
	cookie, err := req.Cookie(config.SessionKey)

	if err != nil {
		return "", err
	}

	if cookie == nil || cookie.Value == "" {
		err := errors.New("empty session cookie")
		return "", err
	}

	return cookie.Value, nil
}
