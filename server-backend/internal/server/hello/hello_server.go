package hello

import (
	"context"

	pb "github.com/1ambda/go-ref/server-backend/pkg/api"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type HelloServer struct{}

func (s *HelloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	actualName, err := GetName(in.Name)
	if err != nil {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		log := logger.Sugar()
		log.Errorw("Failed to retrieve name",
			"error", err)
	}
	return &pb.HelloReply{Message: "Hello " + actualName}, nil
}

func GetName(name string) (string, error) {
	if name == "2ambda" {
		err := errors.New("Invalid name: " + name)
		return "", err
	}

	if name == "1ambda" {
		name = "Kun"
	}

	return name, nil
}
