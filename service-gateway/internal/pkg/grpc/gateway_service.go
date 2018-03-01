package grpcservice

import (
	pb "github.com/1ambda/go-ref/service-gateway/pkg/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/grpclog"
	"time"
	"go.uber.org/zap"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/service"
	"github.com/satori/go.uuid"
)

type GatewayService struct {
	realtimeStatService service.RealtimeStatService
}

func NewGatewayService() *GatewayService {
	return &GatewayService{*service.NewRealtimeStatService(),}
}

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

	stream.SetTrailer(metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-stream"))
	return nil
}

func (s *GatewayService) SubscribeCurrentUserCount(request *pb.EmptyRequest, stream pb.Gateway_SubscribeCurrentUserCountServer) error {
	stream.SendHeader(metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-stream"))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	recv := make(chan int64)
	uuid := uuid.NewV4().String()
	s.realtimeStatService.IncreaseCurrentUserCountSubscriber(uuid, recv)

	done := false
	for !done {
		select {
		case currentUserCount, ok := <- recv:
			if !ok {
				done = true; break
			}

			count := pb.CountResponse{Count: currentUserCount}
			sugar.Info("Sending gRPC response: SubscribeCurrentUserCount")

			err := stream.Send(&count)
			if err != nil {
				sugar.Errorw("Failing to send gRPC response", "error", err)
				done = true; break
			}

		case <- stream.Context().Done():
			// https://github.com/improbable-eng/grpc-web/issues/57
			sugar.Info("Client disconnected")
			done = true; break
		}

	}

	s.realtimeStatService.DecreaseCurrentUserCountSubscriber(uuid)
	stream.SetTrailer(metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-stream"))
	return nil
}
