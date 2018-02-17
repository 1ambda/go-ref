package rest

import (
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations"
	//"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations/todos"
	//"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
)

func ConfigureAPI(api *operations.GatewayAPI) {
	//api.TodosFindTodosHandler = todos.FindTodosHandler(func(params todos.FindTodos) middleware.Responder {
	//	return todos.NewFindTodosOK()
	//})

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Get Todo")
}



