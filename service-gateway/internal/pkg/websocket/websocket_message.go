package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/ws_model"
)

type WebSocketMessage struct {
	content *[]byte // message text
	event   string  // message type
}

func NewConnectionCountMessage(count int) (*WebSocketMessage, error) {
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
		return nil, err
	}

	return &WebSocketMessage{content: &serialized, event: eventType}, nil
}

func NewLeaderNameMessage(leaderName string) (*WebSocketMessage, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateMasterIdentifier
	stringified := fmt.Sprintf("%s", leaderName)

	message := ws_model.WebSocketRealtimeResponse{
		Header: &ws_model.WebSocketResponseHeader{ResponseType: &eventType},
		Body: &ws_model.WebSocketRealtimeResponseBody{
			Value: &stringified,
		},
	}

	serialized, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return &WebSocketMessage{content: &serialized, event: eventType}, nil
}
