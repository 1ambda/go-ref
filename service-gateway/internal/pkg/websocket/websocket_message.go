package websocket

import (
	"fmt"
	"encoding/json"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/ws_model"

	"github.com/1ambda/go-ref/service-gateway/internal/pkg/logger"
)

type WebSocketMessage struct {
	content *[]byte // message text
	event   string   // message type
}

func buildConnectionCountMessage(count int) (*[]byte, error) {
	responseType := ws_model.WebSocketResponseHeaderResponseTypeUpdateConnectionCount
	stringified := fmt.Sprintf("%d", count)

	message := ws_model.WebSocketRealtimeResponse{
		Header: &ws_model.WebSocketResponseHeader{ResponseType: &responseType,},
		Body: &ws_model.WebSocketRealtimeResponseBody{
			Value: &stringified,
		},
	}

	serialized, err := json.Marshal(message)
	if err != nil {
		logger.Errorw("Failed to marshal UPDATE_CURRENT_CONNECTION_COUNT", "error", err)
		return nil, err
	}

	return &serialized, nil
}

func NewConnectionCountWebsocketMessage(count int) (*WebSocketMessage, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateConnectionCount
	stringified := fmt.Sprintf("%d", count)

	message := ws_model.WebSocketRealtimeResponse{
		Header: &ws_model.WebSocketResponseHeader{ResponseType: &eventType,},
		Body: &ws_model.WebSocketRealtimeResponseBody{
			Value: &stringified,
		},
	}

	serialized, err := json.Marshal(message)
	if err != nil {
		logger.Errorw("Failed to marshal UPDATE_CURRENT_CONNECTION_COUNT", "error", err)
		return nil, err
	}

	return &WebSocketMessage{ content: &serialized, event: eventType}, nil
}
