package main

import (
	"net"

	"github.com/1ambda/go-ref/server-backend/internal/pkg/config"
	"github.com/1ambda/go-ref/server-backend/internal/server/hello"
	pb "github.com/1ambda/go-ref/server-backend/pkg/api"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	port := spec.Port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &hello.HelloServer{})
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
