package distributed

import (
	"context"
	"time"

	"github.com/1ambda/go-ref/service-location/internal/config"
	"github.com/coreos/etcd/clientv3"
)

const CampaignInterval = 5 * time.Second
const ElectionTimeout = 10 * time.Second
const EtcdSessionTTL = 120 // second
const EtcdPutGetTimeout = 5 * time.Second

// distributed storage connector.
type Connector interface {
	Publish(message *Message) error
	Stop()
}

type etcdConnectorImpl struct {
	appCtx context.Context
	etcdClient *clientv3.Client

	serverName string
	endpoints  []string
}

func New(appCtx context.Context, endpoints []string, serverName string) (Connector, error) {
	logger := config.GetLogger()

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		etcdClient.Close()
		if err == context.DeadlineExceeded {
			logger.Fatalw("Failed to connect to etcd cluster",
				"cause", "dial timeout", "error", err)
		}

		logger.Fatalw("Failed to connect to etcd cluster",
			"cause", "unknown error", "error", err)
	}

	return &etcdConnectorImpl{
		appCtx:     appCtx,
		etcdClient: etcdClient,
		endpoints:  endpoints,
		serverName: serverName,
	}, nil
}

func (c *etcdConnectorImpl) Publish(message *Message) error {
	panic("implement me")
}

func (c *etcdConnectorImpl) Stop() {
	panic("implement me")
}

type Message struct {
	key   string
	value string
}
