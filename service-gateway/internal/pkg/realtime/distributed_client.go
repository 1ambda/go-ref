package realtime

import (
	"context"
	"time"

	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"go.uber.org/zap"
	"sync"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
)

type DistributedClient interface {
	Stop()
}

const UpdateStatTaskInterval = 3 * time.Second
const LeaderTaskInterval = 5 * time.Second
const CampaignInterval = 5 * time.Second
const ElectionTimeout = 10 * time.Second
const ElectionPath = "/gateway-leader"
const EtcdSessionTTL = 5 // second
const EtcdPutGetTimeout = 5 * time.Second

type etcdDistributedClient struct {
	client     *clientv3.Client
	session    *concurrency.Session
	lock       sync.RWMutex
	leader     string
	serverName string
}

func NewDistributedClient(appCtx context.Context, endpoints []string, serverName string) DistributedClient {
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

	dClient := &etcdDistributedClient{client: etcdClient, leader: "", serverName: serverName,}

	session, err := concurrency.NewSession(etcdClient, concurrency.WithTTL(EtcdSessionTTL))
	if err != nil {
		session.Close()
		etcdClient.Close()
		logger.Fatalw("Failed to get etcd session", "error", err)
	}
	dClient.session = session

	logger.Infow("Got etdc session", "lease_id", session.Lease())

	go dClient.runLeaderTask(appCtx)
	go dClient.runUpdateStatTask(appCtx)
	go dClient.runWatchTask(appCtx)
	go dClient.runElectionCampaign(appCtx)

	return dClient
}


func (d *etcdDistributedClient) put(appCtx context.Context, key string, value string) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	ctx, cancel := context.WithTimeout(context.Background(), EtcdPutGetTimeout)
	_, err := d.client.Put(ctx, key, value)
	defer cancel()

	if err != nil {
		switch err {
		case context.Canceled:
			logger.Errorf("ctx is canceled by another routine", "error", err)
		case context.DeadlineExceeded:
			logger.Errorf("ctx is attached with a deadline is exceeded", "error", err)
		case rpctypes.ErrEmptyKey:
			logger.Errorf("Empty etdc key", "error", err)
		default:
			logger.Errorf("bad cluster endpoints, which are not etcd servers", "error", err)
		}
	}

}

func (d *etcdDistributedClient) runUpdateStatTask(appCtx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	ticker := time.NewTicker(UpdateStatTaskInterval)

	for {
		select {
		case <-appCtx.Done():
			logger.Infow("Stopping follower task goroutine")
			return

		case <-ticker.C:
			logger.Infow("Follower Task")
			// 1. Update this server specific stats (websocket), ...
		}
	}
}

func (d *etcdDistributedClient) runWatchTask(appCtx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	// TODO: make variables for watched values

	kWsConnections := "gateway/stat/ws-connection-count"
	wchWsConnections := d.client.Watch(appCtx, kWsConnections, clientv3.WithPrefix())

	stop := false
	for !stop {
		select {
		case <-appCtx.Done():
			stop = true
			break

		case watchResponse := <-wchWsConnections:
			if watchResponse.Canceled {
				logger.Errorw("etcd watch channel is about to close", "key", kWsConnections)
				stop = true
				break
			}

			logger.Infow("Watch Task")

			// TODO: Update this server specific stats (websocket), ...
		}
	}

	logger.Infow("Stopping watch task goroutine")
}

func (d *etcdDistributedClient) runLeaderTask(appCtx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	ticker := time.NewTicker(LeaderTaskInterval)

	for {
		select {
		case <-appCtx.Done():
			logger.Infow("Stopping leader task goroutine")
			return

		case <-ticker.C:
			d.lock.RLock()
			leader := d.leader
			d.lock.RUnlock()

			if leader != d.serverName {
				// not the leader
				logger.Infow("Skip leader task", "leader", leader, "server_name", d.serverName)
				continue
			}

			logger.Infow("Process leader task", "leader", leader, "server_name", d.serverName)

			// 1. get all values and calculate it
			// and put the summary again
		}
	}
}

func (d *etcdDistributedClient) runElectionCampaign(appCtx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	ticker := time.NewTicker(CampaignInterval)
	election := concurrency.NewElection(d.session, ElectionPath)

	// initial election campaign attempt
	d.campaign(appCtx, election)

	for {
		select {
		case <-appCtx.Done():
			logger.Infow("Stopping election campaign goroutine")
			return

		case <-ticker.C:
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
	defer electionCancelFunc()

	resp, err := election.Leader(electionCtx)

	if err == nil {
		if len(resp.Kvs) == 0 {
			logger.Warnw("Got invalid Leader response from etcd. Will try next time",
				"path", ElectionPath)
			return
		}

		kv := resp.Kvs[0]
		leader := fmt.Sprintf("%s", kv.Value)
		d.lock.Lock()
		d.leader = leader
		d.lock.Unlock()
		return
	}

	if err != concurrency.ErrElectionNoLeader {
		logger.Errorw("Failed to get the leader",
			"path", ElectionPath, "error", err)
		return
	}

	logger.Infow("No leader found. Will campaign", "path", ElectionPath)
	if err := election.Campaign(electionCtx, d.serverName); err != nil {
		logger.Errorw("Failed to take leadership. Will try next time",
			"path", ElectionPath, "error", err)
		return
	}

	logger.Infow("Took leadership", "path", ElectionPath, "server_name", d.serverName)
}

func (d *etcdDistributedClient) Stop() {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

	d.session.Close()
	d.client.Close()

	logger.Info("Closed etcd client connection")
}
