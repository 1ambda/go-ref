package rest

import (
	"time"

	"github.com/1ambda/go-ref/service-gateway/internal/distributed"
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/browser_history"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/geolocation"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/session"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/gorm"
)

func getCode(e *dto.RestError) int {
	return int(e.Code)
}

func Configure(db *gorm.DB, api *rest_api.GatewayRestAPI, dClient distributed.DistributedClient) {

	/**
	 * Access API
	 */
	api.BrowserHistoryAddOneHandler = browser_history.AddOneHandlerFunc(
		func(params browser_history.AddOneParams) middleware.Responder {
			restResp, restErr := addOneBrowserHistory(params, db, dClient)
			if restErr != nil {
				return browser_history.NewFindAllDefault(getCode(restErr)).WithPayload(restErr)
			}

			return browser_history.NewAddOneCreated().WithPayload(restResp)
		})

	api.BrowserHistoryFindAllHandler = browser_history.FindAllHandlerFunc(
		func(params browser_history.FindAllParams) middleware.Responder {
			pagination, rows, restErr := findAllBrowserHistory(params, db)
			if restErr != nil {
				return browser_history.NewFindAllDefault(getCode(restErr)).WithPayload(restErr)
			}

			return browser_history.NewFindAllOK().WithPayload(&dto.BrowserHistoryWithPagination{
				Pagination: pagination, Rows: rows,
			})
		})

	api.BrowserHistoryFindOneHandler = browser_history.FindOneHandlerFunc(
		func(params browser_history.FindOneParams) middleware.Responder {
			restResp, restErr := findOneBrowserHistory(params, db)
			if restErr != nil {
				return browser_history.NewFindOneDefault(getCode(restErr)).WithPayload(restErr)
			}
			return browser_history.NewFindOneOK().WithPayload(restResp)
		})

	api.BrowserHistoryRemoveOneHandler = browser_history.RemoveOneHandlerFunc(
		func(params browser_history.RemoveOneParams) middleware.Responder {
			restErr := removeOneBrowserHistory(params, db)
			if restErr != nil {
				return browser_history.NewRemoveOneDefault(getCode(restErr))
			}
			return browser_history.NewRemoveOneNoContent()
		})

	/**
	 * Session API
	 */
	api.SessionValidateOrGenerateHandler = session.ValidateOrGenerateHandlerFunc(
		func(params session.ValidateOrGenerateParams) middleware.Responder {
			restResp, restErr := validateOrGenerateSession(params, db)
			if restErr != nil {
				return session.NewValidateOrGenerateDefault(getCode(restErr)).WithPayload(restErr)
			}

			return session.NewValidateOrGenerateOK().WithPayload(restResp)
		})

	/**
	 * Geolocation API
	 */
	 api.GeolocationAddHandler = geolocation.AddHandlerFunc(
	 	func(params geolocation.AddParams) middleware.Responder {
	 		restResp, restErr := addOneGeolocationHistory(params, db)

			if restErr != nil {
				return geolocation.NewAddDefault(getCode(restErr)).WithPayload(restErr)
			}

			return geolocation.NewAddCreated().WithPayload(restResp)
		})
}

func buildRestError(err error, errorType string, code int64) *dto.RestError {
	return &dto.RestError{
		Code:      code,
		Message:   err.Error(),
		Type: errorType,
		Timestamp: time.Now().UTC().String(),
	}
}
