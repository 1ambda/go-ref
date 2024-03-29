package config

import (
	"github.com/kelseyhightower/envconfig"
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

const ENV_LOCAL = "LOCAL"
const ENV_TEST = "TEST"
const ENV_DEV = "DEV"
const ENV_PROD = "PROD"

type Specification struct {
	Debug                  bool     `default:"true"`
	Env                    string   `default:"LOCAL"` // `LOCAL`, `TEST`, `DEV`, `PROD`
	Host                   string   `default:"localhost"`
	WebSocketPort          int      `default:"50001"`
	HttpPort               int      `default:"50002"`
	CorsAllowUrl           string   `default:"localhost:3000"`
	EtcdEndpoints          []string `default:"http://127.0.0.1:2379"`
	LocationServerEndpoint string   `default:"localhost:50003"`
	ServerName             string   `default:"0"`
	MysqlHost              string   `default:"localhost"`
	MysqlPort              string   `default:"3306"`
	MysqlUserName          string   `default:"root"`
	MysqlPassword          string   `default:"root"`
	MysqlDatabase          string   `default:"goref"`
}

var Spec Specification

func init() {
	Spec = GetSpecification()
}

func GetSpecification() Specification {
	var s Specification
	err := envconfig.Process("", &s)
	if err != nil {
		panic("Failed to get specification")
	}

	return s
}

func IsTestEnv(spec Specification) bool {
	return spec.Env == ENV_TEST
}

func IsLocalEnv(spec Specification) bool {
	return spec.Env == ENV_LOCAL
}
