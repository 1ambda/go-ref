package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/1ambda/go-ref/service-location/internal/pkg/config"
	"github.com/1ambda/go-ref/service-location/internal/server/hello"
	"github.com/1ambda/go-ref/service-location/pkg/generated/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	spec := config.GetSpecification()

	logger := config.GetLogger()
	logger.Infow("Starting server...",
		"version", config.Version,
		"build_date", config.BuildDate,
		"git_commit", config.GitCommit,
		"git_branch", config.GitBranch,
		"git_state", config.GitState,
		"git_summary", config.GitSummary,
		"env", spec.Env,
		"grpc_port", spec.GrpcPort,
		"debug", spec.Debug,
		"server_name", spec.ServerName,
	)

	port := spec.GrpcPort
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	// register grpc services
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &hello.HelloServer{})
	reflection.Register(s)

	// register shutdown hook
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		logger.Infow("Stopping server...", "server_name", spec.ServerName)
		s.GracefulStop()
	}()

	// start server
	if err := s.Serve(listener); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}
