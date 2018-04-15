package websocket

import (
	"encoding/json"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/ws_model"
)

type Message struct {
	content *[]byte // message text
	event   string  // message type
}

func NewStringValueMessage(eventType string, count string) (*Message, error) {
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

	return &Message{content: &serialized, event: eventType}, nil
}

func NewGatewayNodeCountMessage(count string) (*Message, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateGatewayNodeCount
	return NewStringValueMessage(eventType, count)
}

func NewWebSocketConnectionCountMessage(count string) (*Message, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateWebSocketConnectionCount
	return NewStringValueMessage(eventType, count)
}

func NewGatewayLeaderNodeNameMessage(leaderName string) (*Message, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateGatewayLeaderNodeName
	return NewStringValueMessage(eventType, leaderName)
}

func NewBrowserHistoryCountMessage(count string) (*Message, error) {
	eventType := ws_model.WebSocketResponseHeaderResponseTypeUpdateBrowserHistoryCount
	return NewStringValueMessage(eventType, count)
}
