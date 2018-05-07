package location

import (
	"context"
	"fmt"
	"sync"

	"github.com/1ambda/go-ref/service-location/internal/distributed"
	"github.com/1ambda/go-ref/service-location/pkg/generated/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type country string

type leader string

type server struct {
	lock sync.RWMutex

	leaders map[country]leader

	connector  distributed.Connector
	serverName string
}

func New(srvName string, connector distributed.Connector) (pb.LocationServer, error) {
	svc := &server{
		connector:  connector,
		serverName: srvName,
		leaders: make(map[country]leader),
	}

	return svc, nil
}

func (s *server) updateLeaders(c country, l leader) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.leaders[c] = l
}

func (s *server) Add(ctx context.Context, in *pb.LocationRequest) (*pb.LocationResponse, error) {

	if in.LocationContext == nil || in.LocationContext.Country == "" || in.LocationContext.SessionId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid LocationRequest")
	}

	c := in.LocationContext.Country

	// (TODO): leader cache
	// get leader
	srvName := s.serverName
	l, err := s.connector.GetLeaderOrCampaign(c, srvName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.updateLeaders(country(c), leader(l))

	if srvName != l {
		message := fmt.Sprintf("%s is not owned by %s but %s has ownership", c, srvName, l)
		return nil, status.Error(codes.InvalidArgument, message)
	}

	resp := &pb.LocationResponse{
		LocationContext: in.LocationContext,
	}

	return resp, nil
}
