package rest

import (
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/access"

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
	api.AccessAddOneHandler = access.AddOneHandlerFunc(
		func(params access.AddOneParams) middleware.Responder {
			restResp, restErr := addOneAccess(params, db, dClient)
			if restErr != nil {
				return access.NewFindAllDefault(getCode(restErr)).WithPayload(restErr)
			}

			return access.NewAddOneCreated().WithPayload(restResp)
		})

	api.AccessFindAllHandler = access.FindAllHandlerFunc(
		func(params access.FindAllParams) middleware.Responder {
			pagination, rows, restErr := findAllAccess(params, db)
			if restErr != nil {
				return access.NewFindAllDefault(getCode(restErr)).WithPayload(restErr)
			}

			return access.NewFindAllOK().WithPayload(&rest_model.FindAllOKBody{
				Pagination: pagination, Rows: rows,
			})
		})

	api.AccessFindOneHandler = access.FindOneHandlerFunc(
		func(params access.FindOneParams) middleware.Responder {
			restResp, restErr := findOneAccess(params, db)
			if restErr != nil {
				return access.NewFindOneDefault(getCode(restErr)).WithPayload(restErr)
			}
			return access.NewFindOneOK().WithPayload(restResp)
		})

	api.AccessRemoveOneHandler = access.RemoveOneHandlerFunc(
		func(params access.RemoveOneParams) middleware.Responder {
			restErr := removeOneAccess(params, db)
			if restErr != nil {
				return access.NewRemoveOneDefault(getCode(restErr))
			}
			return access.NewRemoveOneNoContent()
		})

	api.AccessUpdateOneHandler = access.UpdateOneHandlerFunc(
		func(params access.UpdateOneParams) middleware.Responder {
			restResp, restErr := updateOneAccess(params, db)
			if restErr != nil {
				return access.NewAddOneDefault(getCode(restErr)).WithPayload(restErr)
			}
			return access.NewUpdateOneOK().WithPayload(restResp)
		})

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


