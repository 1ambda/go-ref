package main

import (
	"os"
	"fmt"
	"net/http"
	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"github.com/go-openapi/runtime"

	"github.com/rs/cors"

	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/1ambda/go-ref/service-gateway/internal/pkg/config"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/model"

	restapi "github.com/1ambda/go-ref/service-gateway/pkg/api/rest"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/rest"
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/service"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/websocket"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// get config
	spec := config.GetSpecification()
	sugar.Infow("Starting server...",
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
	)

	// setup db connection
	sugar.Info("Connecting to MySQL")
	db, err := connectToMySQL(spec)
	defer db.Close()
	if err != nil {
		sugar.Fatal(err)
	}

	sugar.Info("Auto-migrate MySQL tables")
	db.SingularTable(true)
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.Access{})

	// create services
	realtimeService := service.NewRealtimeStatService(db)
	wsClientManager := websocketservice.NewWebSocketClientManager()

	// configure WS server handlers, middlewares
	sugar.Info("Configure WS server")
	mux := http.NewServeMux()
	websocketservice.Configure(mux, wsClientManager)
	//wsCors := cors.New(cors.Options{
	//	AllowedOrigins: []string{"http://localhost:3000"},
	//	AllowCredentials: true,
	//})
	//wsHandler := wsCors.Handler(mux)
	go func() {
		wsServerPort := fmt.Sprintf(":%d", spec.WebSocketPort)
		if err := http.ListenAndServe(wsServerPort, mux); err != nil {
			sugar.Fatalw("failed starting websocket server", "error", err)
		}
	}()

	// configure REST server
	sugar.Info("Configure REST server")
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		sugar.Fatal(err)
	}
	api := operations.NewGatewayAPI(swaggerSpec)

	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			sugar.Fatal(err)
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
	api.Logger = sugar.Infof

	// configure REST server handlers, middlewares
	sugar.Info("Configure REST handlers")
	restservice.Configure(db, api, realtimeService)
	handler := api.Serve(nil)

	sugar.Info("Configure REST middleware")
	handler = cors.Default().Handler(handler)
	server.SetHandler(handler)

	api.ServerShutdown = func() {
		sugar.Info("Handling shutdown hook")
	}

	if err := server.Serve(); err != nil {
		sugar.Fatal(err)
	}
}

func connectToMySQL(spec config.Specification) (*gorm.DB, error) {
	dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		spec.MysqlUserName, spec.MysqlPassword, spec.MysqlHost, spec.MysqlPort, spec.MysqlDatabase)
	db, err := gorm.Open("mysql", dbConnString)

	return db, err
}
