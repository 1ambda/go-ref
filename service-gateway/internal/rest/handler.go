package rest

import (
	"time"

	"github.com/1ambda/go-ref/service-gateway/internal/distributed"
	"github.com/1ambda/go-ref/service-gateway/internal/location"
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/browser_history"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/geolocation"
	"go.uber.org/zap"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/session"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/gorm"
)

type Controller interface {
	Configure(api *rest_api.GatewayRestAPI)
}

type controllerImpl struct {
	db              *gorm.DB
	dConnector      distributed.Connector
	logger          *zap.SugaredLogger
	locationService location.Service
}

func New(db *gorm.DB, dConnector distributed.Connector, locationService location.Service, logger *zap.SugaredLogger) Controller {
	return &controllerImpl{db: db, dConnector: dConnector, locationService: locationService, logger: logger,}
}

func (ctrl *controllerImpl) Configure(api *rest_api.GatewayRestAPI) {

	/**
	 * Access API
	 */
	api.BrowserHistoryAddOneHandler = browser_history.AddOneHandlerFunc(
		func(params browser_history.AddOneParams) middleware.Responder {
			restResp, restErr := ctrl.addOneBrowserHistory(params)
			if restErr != nil {
				return browser_history.NewFindAllDefault(getCode(restErr)).WithPayload(restErr)
			}

			return browser_history.NewAddOneCreated().WithPayload(restResp)
		})

	api.BrowserHistoryFindAllHandler = browser_history.FindAllHandlerFunc(
		func(params browser_history.FindAllParams) middleware.Responder {
			pagination, rows, restErr := ctrl.findAllBrowserHistory(params)
			if restErr != nil {
				return browser_history.NewFindAllDefault(getCode(restErr)).WithPayload(restErr)
			}

			return browser_history.NewFindAllOK().WithPayload(&dto.BrowserHistoryWithPagination{
				Pagination: pagination, Rows: rows,
			})
		})

	api.BrowserHistoryFindOneHandler = browser_history.FindOneHandlerFunc(
		func(params browser_history.FindOneParams) middleware.Responder {
			restResp, restErr := ctrl.findOneBrowserHistory(params)
			if restErr != nil {
				return browser_history.NewFindOneDefault(getCode(restErr)).WithPayload(restErr)
			}
			return browser_history.NewFindOneOK().WithPayload(restResp)
		})

	api.BrowserHistoryRemoveOneHandler = browser_history.RemoveOneHandlerFunc(
		func(params browser_history.RemoveOneParams) middleware.Responder {
			restErr := ctrl.removeOneBrowserHistory(params)
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
			restResp, restErr := ctrl.validateOrGenerateSession(params)
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
			restResp, restErr := ctrl.addOneGeolocationHistory(params)

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
		Type:      errorType,
		Timestamp: time.Now().UTC().String(),
	}
}

func getCode(e *dto.RestError) int {
	return int(e.Code)
}
