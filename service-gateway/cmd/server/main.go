package main

import (
	"os"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/config"
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest"
	internal "github.com/1ambda/go-ref/service-gateway/internal/server/rest"
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations"

	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"github.com/go-openapi/runtime"
)

func main() {
	spec := config.GetSpecification()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()
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
	internal.ConfigureAPI(api)

	// set middleware
	handler := api.Serve(nil)
	server.SetHandler(handler)

	// set shutdown hook
	api.ServerShutdown = func() {}

	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}

