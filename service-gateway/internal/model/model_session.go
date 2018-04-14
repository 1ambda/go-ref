package model

import (
	"time"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
)

const SessionTable = "session"

type Session struct {
	BaseModel

	SessionID    string    `gorm:"not null;unique"`
	ExpiredAt    time.Time
	Refreshed    bool      `gorm:"not null"`
	RefreshCount int       `gorm:"not null"`
}

func ConvertToSessionDTO(record *Session) *rest_model.SessionResponse {
	// return millis
	updatedAt := record.UpdatedAt.UnixNano() / 1000000
	createdAt := record.CreatedAt.UnixNano() / 1000000
	expiredAt := record.ExpiredAt.UnixNano() / 1000000
	refreshCount := int64(record.RefreshCount)

	return &rest_model.SessionResponse{
		SessionID:    &record.SessionID,
		UpdatedAt:    &updatedAt,
		CreatedAt:    &createdAt,
		ExpiredAt:    &expiredAt,
		RefreshCount: &refreshCount,
		Refreshed:    &record.Refreshed,
	}
}
