package grpcservice

import (
	"golang.org/x/net/context"
	pb "github.com/1ambda/go-ref/service-gateway/pkg/grpc"
)

type HelloService struct{}

func (s *HelloService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello" + req.Name}, nil
}


