package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/jinzhu/gorm"
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

	Session   Session `gorm:"foreignkey:SessionID"`
	SessionID string  `gorm:"column:session_id; type:VARCHAR(255) REFERENCES session(session_id)"`
}

func ConvertFromBrowserHistoryDTO(dto *dto.BrowserHistory) *BrowserHistory {
	record := BrowserHistory{
		BrowserName:     *dto.BrowserName,
		BrowserVersion:  *dto.BrowserVersion,
		OsName:          *dto.OsName,
		OsVersion:       *dto.OsVersion,
		IsMobile:        *dto.IsMobile,
		ClientTimezone:  *dto.ClientTimezone,
		ClientTimestamp: *dto.ClientTimestamp,
		Language:        *dto.Language,
		UserAgent:       *dto.UserAgent,
	}

	return &record
}

func ConvertToBrowserHistoryDTO(record *BrowserHistory) *dto.BrowserHistory {
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
