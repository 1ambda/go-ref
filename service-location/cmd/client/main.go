package main

import (
	"github.com/1ambda/go-ref/service-location/internal/pkg/config"
	"github.com/1ambda/go-ref/service-location/pkg/generated/grpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	spec := config.GetSpecification()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()
	log.Infow("Starting client...",
		"build_date", config.Version,
		"build_date", config.BuildDate,
		"git_commit", config.GitCommit,
		"git_branch", config.GitBranch,
		"git_state", config.GitState,
		"git_summary", config.GitSummary,
		"env", spec.Env,
		"host", spec.Host,
		"port", spec.GrpcPort,
		"debug", spec.Debug,
	)

	address := spec.Host + ":" + spec.GrpcPort
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloClient(conn)

	// Contact the server and print out its response.
	name := "1ambda"
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}

	log.Infof("Greeting: %s", r.Message)
}
