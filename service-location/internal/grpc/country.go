package grpc

import (
	"context"

	"github.com/1ambda/go-ref/service-location/pkg/generated/grpc"
)

type locationServiceImpl struct {

}

func (s *locationServiceImpl) AddSession(ctx context.Context, in *pb.LocationRequest) (*pb.LocationResponse, error) {
	return nil, nil
}

