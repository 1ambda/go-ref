package model

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type BaseModel struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"deletedAt"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

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
}
