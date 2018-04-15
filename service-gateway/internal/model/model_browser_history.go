package model

import (
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const BrowserHistoryTable = "browser_history"

type BrowserHistory struct {
	gorm.Model

	BrowserName     string `gorm:"column:browser_name; not null"`
	BrowserVersion  string `gorm:"column:browser_version; not null"`
	OsName          string `gorm:"column:os_name; not null"`
	OsVersion       string `gorm:"column:os_version; not null"`
	IsMobile        bool   `gorm:"column:is_mobile; not null"`
	ClientTimezone  string `gorm:"column:client_timezone; not null"`
	ClientTimestamp string `gorm:"column:client_timestamp; not null"`
	Language        string `gorm:"column:language; not null"`
	UserAgent       string `gorm:"column:user_agent; not null"`

	// foreign keys
	Session   Session `gorm:"foreignkey:SessionID"`
	SessionID string  `gorm:"column:session_id; type:VARCHAR(255) REFERENCES session(session_id)"`
}

func (record *BrowserHistory) ConvertFromBrowserHistoryDTO(dto *dto.BrowserHistory) *BrowserHistory {
	record.BrowserName = *dto.BrowserName
	record.BrowserVersion = *dto.BrowserVersion
	record.OsName = *dto.OsName
	record.OsVersion = *dto.OsVersion
	record.IsMobile = *dto.IsMobile
	record.ClientTimezone = *dto.ClientTimezone
	record.ClientTimestamp = *dto.ClientTimestamp
	record.Language = *dto.Language
	record.UserAgent = *dto.UserAgent

	return record
}

func (record *BrowserHistory) ConvertToBrowserHistoryDTO() *dto.BrowserHistory {
	dto := dto.BrowserHistory{
		ID:              int64(record.ID),
		BrowserName:     &record.BrowserName,
		BrowserVersion:  &record.BrowserVersion,
		OsName:          &record.OsName,
		OsVersion:       &record.OsVersion,
		IsMobile:        &record.IsMobile,
		ClientTimezone:  &record.ClientTimezone,
		ClientTimestamp: &record.ClientTimestamp,
		Language:        &record.Language,
		UserAgent:       &record.UserAgent,
	}

	return &dto
}
