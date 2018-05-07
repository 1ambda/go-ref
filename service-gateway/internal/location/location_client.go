package location

import (
	"context"
	"fmt"
	"time"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Service interface {
	Add(sessionId string, country string) (*pb.LocationResponse, error)
	Close()
}

const clientTimeout = time.Second * 3

type serviceImpl struct {
	grpcConn *grpc.ClientConn
	ctx      context.Context
}

func New(ctx context.Context, endpoint string) (Service, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		message := fmt.Sprintf("Failed to grpc client conn manager for location server %s", endpoint)
		return nil, errors.Wrap(err, message)
	}

	return &serviceImpl{grpcConn: conn, ctx: ctx,}, nil
}

func (c *serviceImpl) Close() {
	logger := config.GetLogger()

	err := c.grpcConn.Close()
	if err != nil {
		logger.Errorw("Failed to close grpc session for location server", "err", err)
	}
}

func (c *serviceImpl) Add(sessionId string, country string) (*pb.LocationResponse, error) {
	client := pb.NewLocationClient(c.grpcConn)

	ctx, cancel := context.WithTimeout(c.ctx, clientTimeout)
	defer cancel()

	req := &pb.LocationRequest{
		LocationContext: &pb.LocationContext{
			SessionId: sessionId, Country: country,
		},
	}
	resp, err := client.Add(ctx, req)

	return resp, err
}
