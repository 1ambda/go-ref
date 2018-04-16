package model

import (
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/jinzhu/gorm"
)

const GeolocationHistoryTable = "geolocation_history"

type GeolocationHistory struct {
	gorm.Model

	// api
	ApiProvider        string `gorm:"column:api_provider;"`
	ApiLanguage        string `gorm:"column:api_language;"`
	ApiVersion         string `gorm:"column:api_version;"`
	ApiDesiredAccuracy int32  `gorm:"column:api_desired_accuracy;"`

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

func (record *GeolocationHistory) ConvertFromDTO(sessionID string, dto *dto.Geolocation) {
	record.SessionID = sessionID

	record.ApiProvider = dto.APIProvider
	record.ApiLanguage = dto.APILanguage
	record.ApiVersion = dto.APIVersion
	record.ApiDesiredAccuracy = dto.APIDesiredAccuracy

	record.Provider = dto.Provider
	record.Timezone = dto.Timezone
	record.IP = dto.IP
	record.GooglePlaceID = dto.GooglePlaceID

	record.Latitude = dto.Latitude
	record.Longitude = dto.Longitude

	record.FormattedAddress = dto.FormattedAddress
	record.CommonName = dto.CommonName
	record.StreetNumber = dto.StreetNumber
	record.Street = dto.Street

	record.Route = dto.Route
	record.Neighborhood = dto.Neighborhood
	record.Town = dto.Town
	record.City = dto.City
	record.Region = dto.Region
	record.PostalCode = dto.PostalCode

	record.State = dto.State
	record.StateCode = dto.StateCode
	record.Country = dto.Country
	record.CountryCode = dto.CountryCode
}

func (record *GeolocationHistory) ConvertToDTO() *dto.Geolocation {
	return &dto.Geolocation{
		SessionID: record.SessionID,

		APIProvider:        record.ApiProvider,
		APILanguage:        record.ApiLanguage,
		APIVersion:         record.ApiVersion,
		APIDesiredAccuracy: record.ApiDesiredAccuracy,

		Provider:      record.Provider,
		Timezone:      record.Timezone,
		IP:            record.IP,
		GooglePlaceID: record.GooglePlaceID,

		Latitude:  record.Latitude,
		Longitude: record.Longitude,

		FormattedAddress: record.FormattedAddress,
		CommonName:       record.CommonName,
		StreetNumber:     record.StreetNumber,
		Street:           record.Street,

		Route:        record.Route,
		Neighborhood: record.Neighborhood,
		Town:         record.Town,
		City:         record.City,
		Region:       record.Region,
		PostalCode:   record.PostalCode,

		State:       record.State,
		StateCode:   record.StateCode,
		Country:     record.Country,
		CountryCode: record.CountryCode,
	}
}
