package rest

import (
	"github.com/1ambda/go-ref/server-gateway/pkg/api/rest/operations"
	"github.com/1ambda/go-ref/server-gateway/pkg/api/rest/operations/todos"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
)

func ConfigureAPI(api *operations.GatewayAPI) {

	api.TodosGetHandler = todos.GetHandlerFunc(func(params todos.GetParams) middleware.Responder {
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		logger.Info("Get Todo")
		return todos.NewGetOK()
	})
}



