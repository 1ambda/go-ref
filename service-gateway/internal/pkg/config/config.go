package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"
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
	Env           string `default:"LOCAL"`
	Debug         bool   `default:"true"`
	WebSocketPort int    `default:"50001"`
	HttpPort      int    `default:"50002"`
	Host          string `default:"localhost"`
	MysqlHost     string `default:"localhost"`
	MysqlPort     string `default:"3306"`
	MysqlUserName string `default:"root"`
	MysqlPassword string `default:"root"`
	MysqlDatabase string `default:"goref"`
}

func GetSpecification() Specification {
	var s Specification
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
