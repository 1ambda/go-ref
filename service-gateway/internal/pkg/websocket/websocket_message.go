package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/ws_model"
	"go.uber.org/zap"
)

type WebSocketMessage struct {
	content *[]byte // message text
	event   string  // message type
}

func NewConnectionCountWebsocketMessage(count int) (*WebSocketMessage, error) {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateConnectionCount
	stringified := fmt.Sprintf("%d", count)

	message := ws_model.WebSocketRealtimeResponse{
		Header: &ws_model.WebSocketResponseHeader{ResponseType: &eventType},
		Body: &ws_model.WebSocketRealtimeResponseBody{
			Value: &stringified,
		},
	}

	serialized, err := json.Marshal(message)
	if err != nil {
		logger.Errorw("Failed to marshal UPDATE_CURRENT_CONNECTION_COUNT", "error", err)
		return nil, err
	}

	return &WebSocketMessage{content: &serialized, event: eventType}, nil
}
