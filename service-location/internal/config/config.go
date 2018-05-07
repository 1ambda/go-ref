package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

var (
	// These fields are populated by govvv
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	Version    string
)

const ENV_LOCAL = "LOCAL"
const ENV_TEST = "TEST"
const ENV_DEV = "DEV"
const ENV_PROD = "PROD"
const SERVICE_NAME = "service-location"

type Specification struct {
	Debug         bool     `default:"true"`
	Env           string   `default:"LOCAL"`
	EtcdEndpoints []string `default:"http://127.0.0.1:2379"`
	ServerName    string   `default:"0"`
	GrpcPort      string   `default:"50003"`
	Host          string   `default:"localhost"`
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

func GetServerName() string {
	return fmt.Sprintf("%s-%s", SERVICE_NAME, Spec.ServerName)
}

func IsTestEnv(spec Specification) bool {
	return spec.Env == ENV_TEST
}

func IsLocalEnv(spec Specification) bool {
	return spec.Env == ENV_LOCAL
}
