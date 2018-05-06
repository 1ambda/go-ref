package location

import (
	"context"
	"sync"

	"github.com/1ambda/go-ref/service-location/internal/distributed"
	"github.com/1ambda/go-ref/service-location/pkg/generated/grpc"
)

type country string

type leader string

type service struct {
	lock sync.RWMutex

	leaders map[country]leader

	connector distributed.Connector
}

func New(connector distributed.Connector) (pb.LocationServer, error) {
	svc := &service{
		connector: connector,
	}

	return svc, nil
}

func (s *service) Add(ctx context.Context, in *pb.LocationRequest) (*pb.LocationResponse, error) {

	// get leader


	// if no leader, do campaign
	// get leader name eventually

	// (TODO): leader name cache
	// send sessionId to country key
	//

	return nil, nil
}
