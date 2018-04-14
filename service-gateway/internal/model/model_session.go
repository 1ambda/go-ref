package model

import (
	"time"
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/jinzhu/gorm"
)

const SessionTable = "session"

type Session struct {
	gorm.Model

	SessionID    string    `gorm:"column:session_id; not null; unique_index"`
	ExpiredAt    time.Time `gorm:"column:expired_at"`
	Refreshed    bool      `gorm:"column:refreshed; not null"`
	RefreshCount int       `gorm:"column:refresh_count; not null"`
}

func ConvertToSessionDTO(record *Session) *dto.SessionResponse {
	// return millis
	updatedAt := record.UpdatedAt.UnixNano() / 1000000
	createdAt := record.CreatedAt.UnixNano() / 1000000
	expiredAt := record.ExpiredAt.UnixNano() / 1000000
	refreshCount := int64(record.RefreshCount)

	return &dto.SessionResponse{
		SessionID:    &record.SessionID,
		UpdatedAt:    &updatedAt,
		CreatedAt:    &createdAt,
		ExpiredAt:    &expiredAt,
		RefreshCount: &refreshCount,
		Refreshed:    &record.Refreshed,
	}
}
