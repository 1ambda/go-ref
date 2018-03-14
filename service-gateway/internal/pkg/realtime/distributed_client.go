package realtime

import (
	"context"
	"time"

	"fmt"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/config"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"go.uber.org/zap"
	"sync"
)

type DistributedClient interface {
	run(ctx context.Context)
	Stop()
}

const LeaderCheckInterval = 5 * time.Second
const ElectionTimeout = 10 * time.Second
const ElectionPath = "/election-summary"

type etcdDistributedClient struct {
	client  *clientv3.Client
	session *concurrency.Session
	lock    sync.Mutex
}

func NewDistributedClient(appCtx context.Context, endpoints []string) DistributedClient {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		etcdClient.Close()
		if err == context.DeadlineExceeded {
			logger.Fatalw("Failed to connect etcd due to dial timeout", "error", err)
		}

		logger.Fatalw("Failed to connect etcd due to unknown error", "error", err)
	}

	dClient := &etcdDistributedClient{client: etcdClient}

	session, err := concurrency.NewSession(etcdClient)
	if err != nil {
		session.Close()
		etcdClient.Close()
		logger.Fatalw("Failed to get etcd session", "error", err)
	}
	dClient.session = session

	logger.Infow("Got etdc session", "lease_id", session.Lease())

	go dClient.run(appCtx)
	go dClient.runElectionCampaign(appCtx)

	return dClient
}

func (d *etcdDistributedClient) run(appCtx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	logger.Infow("Running distributed client main loop",
		"lease_id", d.session.Lease())

	for {
		select {
		case <-appCtx.Done():
			logger.Infow("Stopping distributed client main loop",
				"lease_id", d.session.Lease())
			return

		}
	}
}

func (d *etcdDistributedClient) runElectionCampaign(appCtx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	leaderCheckTicker := time.NewTicker(LeaderCheckInterval)
	election := concurrency.NewElection(d.session, ElectionPath)

	d.campaign(appCtx, election)

	for {
		select {
		case <-appCtx.Done():
			logger.Infow("Stopping election campaign goroutine",
				"lease_id", d.session.Lease())
			return

		case <-leaderCheckTicker.C:
			d.campaign(appCtx, election)
		}
	}
}

// Check leader and if there is a no leader, try to take leadership.
func (d *etcdDistributedClient) campaign(appCtx context.Context, election *concurrency.Election) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()
	electionCtx, electionCancelFunc := context.WithTimeout(appCtx, ElectionTimeout)

	resp, err := election.Leader(electionCtx)

	if err == nil {
		if len(resp.Kvs) == 0 {
			logger.Warnw("Got invalid Leader response from etcd. Will try next time",
				"path", ElectionPath)
			return
		}

		kv := resp.Kvs[0]
		leader := fmt.Sprintf("%s", kv.Value)
		logger.Infow("Got current leader", "path", ElectionPath, "leader", leader)
		return
	}

	if err != concurrency.ErrElectionNoLeader {
		logger.Errorw("Failed to get the leader",
			"path", ElectionPath, "error", err)
		return
	}

	logger.Infow("No leader found. Will campaign", "path", ElectionPath)
	if err := election.Campaign(electionCtx, config.Spec.ServiceName); err != nil {
		logger.Errorw("Failed to take leadership. Will try next time",
			"path", ElectionPath, "error", err)
		return
	}

	logger.Infow("Took leadership", "path", ElectionPath, "server_name", config.Spec.ServiceName)

	electionCancelFunc()
}

func (d *etcdDistributedClient) Stop() {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

	d.session.Close()
	d.client.Close()

	logger.Info("Closed etcd client connection")
}
