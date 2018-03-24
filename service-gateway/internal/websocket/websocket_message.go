package websocket

import (
	"encoding/json"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/ws_model"
)

type WebSocketMessage struct {
	content *[]byte // message text
	event   string  // message type
}

func NewStringWebSocketMessage(eventType string, count string) (*WebSocketMessage, error) {
	message := ws_model.WebSocketRealtimeResponse{
		Header: &ws_model.WebSocketResponseHeader{ResponseType: &eventType},
		Body: &ws_model.WebSocketRealtimeResponseBody{
			Value: &count,
		},
	}

	serialized, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return &WebSocketMessage{content: &serialized, event: eventType}, nil
}

func NewNodeCountMessage(count string) (*WebSocketMessage, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateNodeCount
	return NewStringWebSocketMessage(eventType, count)
}

func NewConnectionCountMessage(count string) (*WebSocketMessage, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateConnectionCount
	return NewStringWebSocketMessage(eventType, count)
}

func NewLeaderNameMessage(leaderName string) (*WebSocketMessage, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateMasterIdentifier
	return NewStringWebSocketMessage(eventType, leaderName)
}

func NewTotalAccessCountMessage(count string) (*WebSocketMessage, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateTotalAccessCount
	return NewStringWebSocketMessage(eventType, count)
}
