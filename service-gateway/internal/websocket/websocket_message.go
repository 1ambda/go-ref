package websocket

import (
	"encoding/json"
	"time"

	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/ws_model"
	"github.com/go-openapi/swag"
)

type Message struct {
	content *[]byte // message text
	event   string  // message type
}

func NewStringValueMessage(eventType string, count string) (*Message, error) {
	message := dto.WebSocketRealtimeResponse{
		Header: &dto.WebSocketResponseHeader{ResponseType: &eventType},
		Body: &dto.WebSocketRealtimeResponseBody{
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
	eventType := dto.WebSocketResponseHeaderResponseTypeUpdateGatewayNodeCount
	return NewStringValueMessage(eventType, count)
}

func NewWebSocketConnectionCountMessage(count string) (*Message, error) {
	eventType := dto.WebSocketResponseHeaderResponseTypeUpdateWebSocketConnectionCount
	return NewStringValueMessage(eventType, count)
}

func NewGatewayLeaderNodeNameMessage(leaderName string) (*Message, error) {
	eventType := dto.WebSocketResponseHeaderResponseTypeUpdateGatewayLeaderNodeName
	return NewStringValueMessage(eventType, leaderName)
}

func NewBrowserHistoryCountMessage(count string) (*Message, error) {
	eventType := dto.WebSocketResponseHeaderResponseTypeUpdateBrowserHistoryCount
	return NewStringValueMessage(eventType, count)
}

func NewErrorMessage(err error, code int64) (*Message, error) {
	eventType := dto.WebSocketResponseHeaderResponseTypeError
	wsErr := dto.WebSocketError{
		Code:      code,
		Message:   swag.String(err.Error()),
		Timestamp: time.Now().UTC().String(),
	}

	message := dto.WebSocketRealtimeResponse{
		Header: &dto.WebSocketResponseHeader{ResponseType: &eventType, Error: &wsErr,},
		Body:   &dto.WebSocketRealtimeResponseBody{},
	}

	serialized, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return &Message{content: &serialized, event: eventType}, nil
}
