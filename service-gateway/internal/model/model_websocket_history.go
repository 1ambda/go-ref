package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

const WebsocketHistoryTable = "websocket_history"

type WebsocketHistory struct {
	gorm.Model

	websocketSessionID string `gorm:"column:websocket_session_id; not null"`
	ServerID           string `gorm:"column:server_id; not null"`
	ServerVersion      string `gorm:"column:server_version; not null"`

	ClosedAt    *time.Time `gorm:"column:closed_at;"`
	CloseReason string     `gorm:"column:close_reason;"`

	Session   Session `gorm:"foreignkey:SessionID"`
	SessionID string  `gorm:"column:session_id; type:VARCHAR(255) REFERENCES session(session_id)"`
}
