package realtime

import (
	"time"
	"context"

	"go.uber.org/zap"
	"github.com/coreos/etcd/clientv3"
)

type DistributedClient interface {
	Close()
}

type etcdDistributedClientImpl struct {
	client *clientv3.Client
}

func NewDistributedClient(endpoints []string) *etcdDistributedClientImpl {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		if err == context.DeadlineExceeded {
			logger.Fatalw("Failed to connect etcd due to dial timeout", "error", err)
		}

		logger.Fatalw("Failed to connect etcd due to unknown error", "error", err)
	}

	return &etcdDistributedClientImpl{client: cli,}
}

func (d *etcdDistributedClientImpl) Stop() {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()


	d.client.Close()

	logger.Info("Closed etcd client connection")
}
