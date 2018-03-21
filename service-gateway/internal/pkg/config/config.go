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
	Env           string   `default:"LOCAL"`
	EtcdEndpoints []string `default:"http://127.0.0.1:2379"`
	ServerName    string   `default:"gateway-0"`
	Debug         bool     `default:"true"`
	WebSocketPort int      `default:"50001"`
	HttpPort      int      `default:"50002"`
	Host          string   `default:"localhost"`
	MysqlHost     string   `default:"localhost"`
	MysqlPort     string   `default:"3306"`
	MysqlUserName string   `default:"root"`
	MysqlPassword string   `default:"root"`
	MysqlDatabase string   `default:"goref"`
}

var Spec Specification

func init() {
	Spec = GetSpecification()

	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

	logger.Infow("Starting server...",
		"version", Version,
		"build_date", BuildDate,
		"git_commit", GitCommit,
		"git_branch", GitBranch,
		"git_state", GitState,
		"git_summary", GitSummary,
		"env", Spec.Env,
		"websocket_port", Spec.WebSocketPort,
		"http_port", Spec.HttpPort,
		"debug", Spec.Debug,
		"etcd_endpoints", Spec.EtcdEndpoints,
		"server_name", Spec.ServerName,
	)
}

func GetSpecification() Specification {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	var s Specification
	err := envconfig.Process("", &s)
	if err != nil {
		logger.Fatalw("Failed to create startup specification", "error", err)
	}

	return s
}
