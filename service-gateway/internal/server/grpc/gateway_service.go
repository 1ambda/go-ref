package grpcservice

import (
	pb "github.com/1ambda/go-ref/service-gateway/pkg/grpc"
	"google.golang.org/grpc/metadata"
)

type GatewayService struct{}

func (s *GatewayService) SubscribeTotalAccessCount(request *pb.Empty, stream pb.Gateway_SubscribeTotalAccessCountServer) error {
	stream.SendHeader(metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-stream"))

	count := pb.Count{Count: 3}

	if err := stream.Send(&count); err != nil {
		return err
	}

	stream.SetTrailer(metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-stream"))
	return nil
}

func (s *GatewayService) SubscribeCurrentUserCount(request *pb.Empty, stream pb.Gateway_SubscribeCurrentUserCountServer) error {
	stream.SendHeader(metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-stream"))

	count := pb.Count{Count: 3}

	if err := stream.Send(&count); err != nil {
		return err
	}

	stream.SetTrailer(metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-stream"))
	return nil
}


