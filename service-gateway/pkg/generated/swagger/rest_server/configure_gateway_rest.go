// This file is safe to edit. Once it exists it will not be overwritten

package rest_server

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/access"
)

//go:generate swagger generate server --target ../pkg/generated/swagger --name  --spec ../../schema/swagger/gateway-rest.yml --api-package rest_api --model-package rest_model --server-package rest_server --exclude-main

func configureFlags(api *rest_api.GatewayRestAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *rest_api.GatewayRestAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.AccessAddOneHandler = access.AddOneHandlerFunc(func(params access.AddOneParams) middleware.Responder {
		return middleware.NotImplemented("operation access.AddOne has not yet been implemented")
	})
	api.AccessFindAllHandler = access.FindAllHandlerFunc(func(params access.FindAllParams) middleware.Responder {
		return middleware.NotImplemented("operation access.FindAll has not yet been implemented")
	})
	api.AccessFindOneHandler = access.FindOneHandlerFunc(func(params access.FindOneParams) middleware.Responder {
		return middleware.NotImplemented("operation access.FindOne has not yet been implemented")
	})
	api.AccessRemoveOneHandler = access.RemoveOneHandlerFunc(func(params access.RemoveOneParams) middleware.Responder {
		return middleware.NotImplemented("operation access.RemoveOne has not yet been implemented")
	})
	api.AccessUpdateOneHandler = access.UpdateOneHandlerFunc(func(params access.UpdateOneParams) middleware.Responder {
		return middleware.NotImplemented("operation access.UpdateOne has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
