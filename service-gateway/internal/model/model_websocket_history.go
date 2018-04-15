package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

const WebsocketHistoryTable = "websocket_history"

type WebsocketHistory struct {
	gorm.Model

	ClosedAt    *time.Time `gorm:"column:closed_at;"`
	CloseReason string     `gorm:"column:close_reason;"`

	Session     Session `gorm:"foreignkey:SessionID"`
	SessionID   string  `gorm:"column:session_id; type:VARCHAR(255) REFERENCES session(session_id)"`
	WebsocketID string  `gorm:"column:websocket_id; not null"`
}

func (record *WebsocketHistory) NewWebSocketHistory(sessionID string, websocketID string) {
	record.SessionID = sessionID
	record.WebsocketID = websocketID
}
