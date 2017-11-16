package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

var (
	// These fields are populated by govvv
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
	Version    string
)

type Specification struct {
	Env        string `envconfig:"ENV",default:"LOCAL"`
	Debug      bool   `envconfig:"DEBUG",default:"false"`
	ServerPort string `envconfig:"PORT",default:"50001"`
	ServerHost string `envconfig:"HOST",default:"localhost"`
}

func GetSpecification() Specification {
	var s Specification
	err := envconfig.Process("", &s)
	if err != nil {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		log := logger.Sugar()
		log.Fatal(err.Error())
	}

	return s
}
