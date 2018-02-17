package rest

import (
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations"
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations/access"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
)

func ConfigureAPI(api *operations.GatewayAPI) {
	api.AccessAddOneHandler = access.AddOneHandlerFunc(
		func(params access.AddOneParams) middleware.Responder {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			sugar := logger.Sugar()
			sugar.Infow("Creating Access",
				"access", params.Body,
			)

			return access.NewAddOneCreated().WithPayload(params.Body)
		})

}
