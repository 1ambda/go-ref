package grpcservice

import (
	pb "github.com/1ambda/go-ref/service-gateway/pkg/grpc"
	"google.golang.org/grpc/metadata"
	"go.uber.org/zap"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/service"
	"github.com/satori/go.uuid"
)

type GatewayService struct {
	RealtimeStatService *service.RealtimeStatService
}

func NewGatewayService(r *service.RealtimeStatService) *GatewayService {
	return &GatewayService{r}
}



func (s *GatewayService) SubscribeTotalAccessCount(request *pb.EmptyRequest, stream pb.Gateway_SubscribeTotalAccessCountServer) error {
	stream.SendHeader(metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-stream"))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	recv := make(chan int64)
	uuid := uuid.NewV4().String()
	s.RealtimeStatService.AddTotalAccessCountSubscriber(uuid, recv)
	s.RealtimeStatService.BroadcastToTalAccessCount()

	done := false
	for !done {
		select {
		case totalAccessUserCount, ok := <- recv:
			if !ok {
				done = true; break
			}

			count := pb.CountResponse{Count: totalAccessUserCount}
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

	s.RealtimeStatService.DeleteTotalAccessCountSubscriber(uuid)
	return nil
}

func (s *GatewayService) SubscribeCurrentUserCount(request *pb.EmptyRequest, stream pb.Gateway_SubscribeCurrentUserCountServer) error {
	stream.SendHeader(metadata.Pairs("Pre-Response-Metadata", "Is-sent-as-headers-stream"))

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	recv := make(chan int64)
	uuid := uuid.NewV4().String()
	s.RealtimeStatService.IncreaseCurrentUserCountSubscriber(uuid, recv)

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

	s.RealtimeStatService.DecreaseCurrentUserCountSubscriber(uuid)
	stream.SetTrailer(metadata.Pairs("Post-Response-Metadata", "Is-sent-as-trailers-stream"))
	return nil
}
