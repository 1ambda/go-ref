package main

import (
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/jessevdk/go-flags"
	"net/http"
	"os"

	"github.com/rs/cors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/1ambda/go-ref/service-gateway/internal/pkg/config"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/model"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/rest"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/websocket"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api"
	"go.uber.org/zap"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/realtime"
	"context"
)

func main() {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

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
		"debug", spec.Debug,
		"etcd_endpoints", spec.EtcdEndpoints,
		"server_name", "gateway-0",
	)

	// setup db connection
	logger.Info("Connecting to MySQL")
	db, err := connectToMySQL(spec)
	defer db.Close()
	if err != nil {
		logger.Fatalw("Failed to connect MySQL", "error", err)
	}

	logger.Info("Auto-migrate MySQL tables")
	db.SingularTable(true)
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.Access{})

	//
	appCtx, appCancelFunc := context.WithCancel(context.Background())

	// setup etcd
	logger.Info("Configure distributed client (etcd)")
	dClient := realtime.NewDistributedClient(appCtx, spec.EtcdEndpoints)

	// configure WS server handlers, middlewares
	logger.Info("Configure WS server")
	mux := http.NewServeMux()

	wsManager := websocket.Configure(appCtx, mux)
	//wsCors := cors.New(cors.Options{
	//	AllowedOrigins: []string{"http://localhost:3000"},
	//	AllowCredentials: true,
	//})
	//wsHandler := wsCors.Handler(mux)

	go func() {
		wsServerPort := fmt.Sprintf(":%d", spec.WebSocketPort)
		if err := http.ListenAndServe(wsServerPort, mux); err != nil {
			logger.Fatalw("failed starting websocket server", "error", err)
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
	rest.Configure(db, api)
	handler := api.Serve(nil)

	logger.Info("Configure REST middleware")
	handler = cors.Default().Handler(handler)
	server.SetHandler(handler)

	api.ServerShutdown = func() {
		logger.Info("Handling shutdown hook")
		appCancelFunc()
		dClient.Stop()
		<-wsManager.Stop()
	}

	if err := server.Serve(); err != nil {
		logger.Fatalw("Failed to start REST server", "error", err)
	}
}

func connectToMySQL(spec config.Specification) (*gorm.DB, error) {
	dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		spec.MysqlUserName, spec.MysqlPassword, spec.MysqlHost, spec.MysqlPort, spec.MysqlDatabase)
	db, err := gorm.Open("mysql", dbConnString)

	return db, err
}
