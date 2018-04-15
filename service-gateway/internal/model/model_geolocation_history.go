package model

import (
	"github.com/jinzhu/gorm"
)

const GeolocationHistoryTable = "geolocation_history"

type GeolocationHistory struct {
	gorm.Model

	// general
	Provider      string `gorm:"column:provider;"`
	Timezone      string `gorm:"column:timezone;"`
	IP            string `gorm:"column:ip;"`
	GooglePlaceID string `gorm:"column:google_place_id;"`

	// coords
	Latitude  float32 `gorm:"column:latitude;"`
	Longitude float32 `gorm:"column:longitude;"`

	// address
	FormattedAddress string `gorm:"column:formatted_address;"`
	CommonName       string `gorm:"column:common_name;"`
	StreetNumber     string `gorm:"column:street_number;"`
	Street           string `gorm:"column:street;"`
	Route            string `gorm:"column:route;"`
	Neighborhood     string `gorm:"column:neighborhood;"`
	Town             string `gorm:"column:town;"`
	City             string `gorm:"column:city;"`
	Region           string `gorm:"column:region;"`
	PostalCode       string `gorm:"column:postal_code;"`
	State            string `gorm:"column:state;"`
	StateCode        string `gorm:"column:state_code;"`
	Country          string `gorm:"column:country;"`
	CountryCode      string `gorm:"column:country_code;"`

	// foreign keys
	Session   Session `gorm:"foreignkey:SessionID"`
	SessionID string  `gorm:"column:session_id; type:VARCHAR(255) REFERENCES session(session_id)"`
}
