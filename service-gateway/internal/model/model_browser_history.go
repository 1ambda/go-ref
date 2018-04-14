package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/satori/go.uuid"
)

const BrowserHistoryTable = "browser_history"

type BrowserHistory struct {
	BaseModel

	BrowserName    string `gorm:"not null"`
	BrowserVersion string `gorm:"not null"`
	OsName         string `gorm:"not null"`
	OsVersion      string `gorm:"not null"`
	IsMobile       string `gorm:"not null"`
	Timezone       string `gorm:"not null"`
	Timestamp      string `gorm:"not null"`
	Language       string `gorm:"not null"`
	UserAgent      string `gorm:"not null"`
	UUID           string `gorm:"not null"`
}

func ConvertFromBrowserHistoryDTO(dto *rest_model.BrowserHistory) *BrowserHistory {
	uuid := uuid.NewV4()

	record := BrowserHistory{
		BrowserName:    *dto.BrowserName,
		BrowserVersion: *dto.BrowserVersion,
		OsName:         *dto.OsName,
		OsVersion:      *dto.OsVersion,
		IsMobile:       *dto.IsMobile,
		Timezone:       *dto.Timezone,
		Timestamp:      *dto.Timestamp,
		Language:       *dto.Language,
		UserAgent:      *dto.UserAgent,
		UUID:           uuid.String(),
	}

	return &record
}

func ConvertToBrowserHistoryDTO(record *BrowserHistory) *rest_model.BrowserHistory {
	dto := rest_model.BrowserHistory{
		ID:             int64(record.Id),
		BrowserName:    &record.BrowserName,
		BrowserVersion: &record.BrowserVersion,
		OsName:         &record.OsName,
		OsVersion:      &record.OsVersion,
		IsMobile:       &record.IsMobile,
		Timezone:       &record.Timezone,
		Timestamp:      &record.Timestamp,
		Language:       &record.Language,
		UserAgent:      &record.UserAgent,
		UUID:           record.UUID,
	}

	return &dto
}
