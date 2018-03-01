package grpcservice

import (
	pb "github.com/1ambda/go-ref/service-gateway/pkg/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/grpclog"
	"time"
	"go.uber.org/zap"
)

type GatewayService struct{}

func (s *GatewayService) SubscribeTotalAccessCount(request *pb.EmptyRequest, stream pb.Gateway_SubscribeTotalAccessCountServer) error {
	stream.SendHeader(metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-stream"))

	var currentUserCount int64 = 3

	for {
		count := pb.CountResponse{Count: currentUserCount}

		err := stream.Send(&count)
		if err != nil {
			grpclog.Error("got error", err, "from", stream)
			break
		}

		// https://github.com/improbable-eng/grpc-web/issues/57
		err = stream.Context().Err()
		if err != nil {
			break
		}

		currentUserCount += 1
		time.Sleep(time.Second)
	}

	//stream.SetTrailer(metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-stream"))
	return nil
}

func (s *GatewayService) SubscribeCurrentUserCount(request *pb.EmptyRequest, stream pb.Gateway_SubscribeCurrentUserCountServer) error {
	stream.SendHeader(metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-stream"))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	var currentUserCount int64 = 3

	for {
		count := pb.CountResponse{Count: currentUserCount}

		sugar.Infow("Sending gRPC response: SubscribeCurrentUserCount",
			"client", stream.Context())
		err := stream.Send(&count)
		if err != nil {
			sugar.Errorw("Failing to send gRPC response", "error", err)
			break
		}

		// https://github.com/improbable-eng/grpc-web/issues/57
		err = stream.Context().Err()
		if err != nil {
			sugar.Errorw("Client disconnected", "error", err)
			break
		}

		currentUserCount += 1
		time.Sleep(time.Second)
	}

	stream.SetTrailer(metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-stream"))
	return nil
}


