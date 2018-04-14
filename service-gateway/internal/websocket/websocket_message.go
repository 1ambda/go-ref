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

func NewGatewayNodeCountMessage(count string) (*WebSocketMessage, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateGatewayNodeCount
	return NewStringWebSocketMessage(eventType, count)
}

func NewWebSocketConnectionCountMessage(count string) (*WebSocketMessage, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateWebSocketConnectionCount
	return NewStringWebSocketMessage(eventType, count)
}

func NewGatewayLeaderNodeNameMessage(leaderName string) (*WebSocketMessage, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateGatewayLeaderNodeName
	return NewStringWebSocketMessage(eventType, leaderName)
}

func NewBrowserHistoryCountMessage(count string) (*WebSocketMessage, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateBrowserHistoryCount
	return NewStringWebSocketMessage(eventType, count)
}
