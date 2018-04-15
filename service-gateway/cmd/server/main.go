package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/jessevdk/go-flags"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
	"github.com/1ambda/go-ref/service-gateway/internal/distributed"
	"github.com/1ambda/go-ref/service-gateway/internal/model"
	"github.com/1ambda/go-ref/service-gateway/internal/rest"
	"github.com/1ambda/go-ref/service-gateway/internal/websocket"
	"github.com/rs/cors"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api"
)

func main() {
	logger := config.GetLogger()

	// get config
	spec := config.Spec
	logger.Infow("Starting server...",
		"version", config.Version,
		"build_date", config.BuildDate,
		"git_commit", config.GitCommit,
		"git_branch", config.GitBranch,
		"git_state", config.GitState,
		"git_summary", config.GitSummary,
		"env", spec.Env,
		"websocket_port", spec.WebSocketPort,
		"http_port", spec.HttpPort,
		"cors_allow_Url", spec.CorsAllowUrl,
		"debug", spec.Debug,
		"etcd_endpoints", spec.EtcdEndpoints,
		"server_name", spec.ServerName,
	)

	// setup db connection
	logger.Info("Configure database")
	db := model.GetDatabase(spec)
	defer db.Close()

	// create app context
	appCtx, appCancelFunc := context.WithCancel(context.Background())

	// configure WS server handlers, middlewares
	logger.Info("Configure WS server")
	mux := http.NewServeMux()
	wsManager := websocket.Configure(appCtx, mux)

	// setup etcd client
	logger.Info("Configure distributed client (etcd)")
	dClient := distributed.NewDistributedClient(appCtx, spec.EtcdEndpoints, spec.ServerName, wsManager)

	wsServerPort := fmt.Sprintf(":%d", spec.WebSocketPort)
	logger.Infof("Serving gateway ws at http://127.0.0.1:%d", spec.WebSocketPort)
	wsServer := &http.Server{Addr: wsServerPort, Handler: mux}

	go func() {
		if err := wsServer.ListenAndServe(); err != nil {
			logger.Infof("Stopped serving gateway ws at http://127.0.0.1:%d", spec.WebSocketPort)
		}
	}()

	// configure REST server
	logger.Info("Configure REST server")
	swaggerSpec, err := loads.Analyzed(rest_server.SwaggerJSON, "")
	if err != nil {
		logger.Fatalw("Failed to configure REST server", "error", err)
	}
	api := rest_api.NewGatewayRestAPI(swaggerSpec)

	server := rest_server.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			logger.Fatalw("Failed to parse command-line option for REST server", "error", err)
		}
	}
	server.Host = spec.Host
	server.Port = spec.HttpPort
	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.Logger = logger.Infof

	// configure REST server handlers, middlewares
	logger.Info("Configure REST handlers")
	rest.Configure(db, api, dClient)
	handler := api.Serve(nil)

	logger.Info("Configure REST middleware")
	handler = cors.AllowAll().Handler(handler)
	server.SetHandler(handler)

	api.ServerShutdown = func() {
		logger.Info("Handling shutdown hook")
		appCancelFunc()
		dClient.Stop()

		<-wsManager.Stop()
		if err := wsServer.Shutdown(nil); err != nil {
			logger.Errorw("Failed to shutdown wsServer gracefully", "error", err)
		}
	}

	if err := server.Serve(); err != nil {
		logger.Fatalw("Failed to start REST server", "error", err)
	}
}
