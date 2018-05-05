package distributed

import (
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
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
	appCtx     context.Context
	etcdClient *clientv3.Client

	serverName string
	endpoints  []string
}

func New(appCtx context.Context, endpoints []string, serverName string) (Connector, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		if err == context.DeadlineExceeded {
			err = errors.Wrap(err, "Failed to connect to etcd cluster (DialTimeout)")
			return nil, err
		}

		err = errors.Wrap(err, "Failed to connect to etcd cluster (Unknown)")
		return nil, err
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

func (c *etcdConnectorImpl) run() {

}

type Message struct {
	key   string
	value string
}
