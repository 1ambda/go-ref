package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type BaseModel struct {
	Id        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

var AccessTable = "access"

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
