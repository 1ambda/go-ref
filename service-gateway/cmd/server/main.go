package main

import (
	"os"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/config"
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest"
	internal "github.com/1ambda/go-ref/service-gateway/internal/server/rest"
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations"

	"github.com/rs/cors"

	"github.com/go-openapi/loads"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"github.com/go-openapi/runtime"

	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/model"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()

	// get config
	spec := config.GetSpecification()
	log.Infow("Starting server...",
		"version", config.Version,
		"build_date", config.BuildDate,
		"git_commit", config.GitCommit,
		"git_branch", config.GitBranch,
		"git_state", config.GitState,
		"git_summary", config.GitSummary,
		"env", spec.Env,
		"port", spec.Port,
		"debug", spec.Debug,
	)

	// setup db connection
	dbConnString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		spec.MysqlUserName, spec.MysqlPassword, spec.MysqlHost, spec.MysqlPort, spec.MysqlDatabase)
	db, err := gorm.Open("mysql", dbConnString)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	db.SingularTable(true)
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.Access{})

	// setup API
	swaggerSpec, err := loads.Analyzed(rest.SwaggerJSON, "")
	if err != nil {
		log.Fatal(err)
	}

	api := operations.NewGatewayAPI(swaggerSpec)

	server := rest.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Gateway API"
	parser.LongDescription = "API Spec for Gateway"

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatal(err)
		}
	}
	server.Host = spec.Host
	server.Port = spec.Port

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
	api.Logger = log.Infof
	internal.ConfigureAPI(db, api)

	// set middlewares
	handler := api.Serve(nil)
	handler = cors.Default().Handler(handler)
	server.SetHandler(handler)

	// set shutdown hook
	api.ServerShutdown = func() {}

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
