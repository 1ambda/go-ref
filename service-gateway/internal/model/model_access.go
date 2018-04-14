package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/satori/go.uuid"
)

const AccessTable = "access"

type Access struct {
	BaseModel

	BrowserName    string
	BrowserVersion string
	OsName         string
	OsVersion      string
	IsMobile       string
	Timezone       string
	Timestamp      string
	Language       string
	UserAgent      string
	UUID           string // uuid v4
}

func ConvertFromAccessDTO(swaggerModel *rest_model.Access) *Access {
	uuid := uuid.NewV4()

	record := Access{
		BrowserName:    *swaggerModel.BrowserName,
		BrowserVersion: *swaggerModel.BrowserVersion,
		OsName:         *swaggerModel.OsName,
		OsVersion:      *swaggerModel.OsVersion,
		IsMobile:       *swaggerModel.IsMobile,
		Timezone:       *swaggerModel.Timezone,
		Timestamp:      *swaggerModel.Timestamp,
		Language:       *swaggerModel.Language,
		UserAgent:      *swaggerModel.UserAgent,
		UUID:           uuid.String(),
	}

	return &record
}

func ConvertToAccessDTO(record *Access) *rest_model.Access {
	swaggerModel := rest_model.Access{
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

	return &swaggerModel
}
