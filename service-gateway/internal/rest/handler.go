package rest

import (
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/browser_history"

	"github.com/1ambda/go-ref/service-gateway/internal/distributed"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/session"
)

func getCode(e *rest_model.Error) int {
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

			return browser_history.NewFindAllOK().WithPayload(&rest_model.FindAllOKBody{
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

	// TOOD(1ambda): PUT doens't work :(
	//api.AccessUpdateOneHandler = browser_history.UpdateOneHandlerFunc(
	//	func(params browser_history.UpdateOneParams) middleware.Responder {
	//		restResp, restErr := updateOneAccess(params, db)
	//		if restErr != nil {
	//			return browser_history.NewAddOneDefault(getCode(restErr)).WithPayload(restErr)
	//		}
	//		return browser_history.NewUpdateOneOK().WithPayload(restResp)
	//	})

	/**
	* Session API
	*/
	api.SessionValidateOrGenerateHandler = session.ValidateOrGenerateHandlerFunc(
		func(params session.ValidateOrGenerateParams) middleware.Responder {
			restResp, restErr := validateOrGenerateSession(params, db, dClient)
			if restErr != nil {
				return session.NewValidateOrGenerateDefault(getCode(restErr)).WithPayload(restErr)
			}

			return session.NewValidateOrGenerateOK().WithPayload(restResp)
		})
}

func buildRestError(err error, code int64) *rest_model.Error {
	return &rest_model.Error{
		Code:      code,
		Message:   swag.String(err.Error()),
		Timestamp: time.Now().UTC().String(),
	}
}
