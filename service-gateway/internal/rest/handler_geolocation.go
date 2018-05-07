package rest

import (
	"github.com/1ambda/go-ref/service-gateway/internal/model"
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/geolocation"
)

func (ctrl *controllerImpl) addOneGeolocationHistory(params geolocation.AddParams) (*dto.Geolocation, *dto.RestError) {
	logger := ctrl.logger
	db := ctrl.db
	locationSvc := ctrl.locationService

	sessionID, restErr := getSessionCookieForRest(params.HTTPRequest, db)
	if restErr != nil {
		return nil, restErr
	}

	logger.Infow("Creating Geolocation record",
		"request", params.Body, "session", sessionID)

	record := &model.GeolocationHistory{}
	record.ConvertFromDTO(sessionID, params.Body)

	if err := db.Create(record).Error; err != nil {
		logger.Errorw("Failed to create new Geolocation record",
			"session_id", sessionID, "error", err)
		restError := buildRestError(err, dto.RestErrorTypeInternalServer, 500)
		return nil, restError
	}

	_, err := locationSvc.Add(sessionID, record.Country)
	if err != nil {
		logger.Errorw("Failed to send geolocation info to location server", "error", err)
		return nil, buildRestError(err, dto.RestErrorTypeInternalServer, 500)
	}

	return record.ConvertToDTO(), nil
}
