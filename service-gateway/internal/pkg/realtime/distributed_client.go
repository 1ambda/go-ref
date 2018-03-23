package realtime

import (
	"context"
	"time"

	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"go.uber.org/zap"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/websocket"
)

type DistributedClient interface {
	Publish(message *DistributedMessage)
	Stop()
}

const CampaignInterval = 5 * time.Second
const ElectionTimeout = 10 * time.Second
const ElectionPath = "/gateway-leader"
const EtcdSessionTTL = 120 // second
const EtcdPutGetTimeout = 5 * time.Second

type etcdDistributedClient struct {
	client     *clientv3.Client
	session    *concurrency.Session
	leader     string
	serverName string

	publishChan chan *DistributedMessage

	wsManager  websocket.WebSocketManager
}

func NewDistributedClient(appCtx context.Context, endpoints []string,
	serverName string, wsManager websocket.WebSocketManager) DistributedClient {

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

	dClient := &etcdDistributedClient{
		client:      etcdClient, leader: "", serverName: serverName, wsManager: wsManager,
		publishChan: make(chan *DistributedMessage),
	}

	session, err := concurrency.NewSession(etcdClient, concurrency.WithTTL(EtcdSessionTTL))
	if err != nil {
		session.Close()
		etcdClient.Close()
		logger.Fatalw("Failed to get etcd session", "error", err)
	}
	dClient.session = session

	logger.Infow("Got etdc session", "lease_id", session.Lease())

	go dClient.runPublishTask(appCtx)
	go dClient.runWatchTask(appCtx)
	go dClient.runElectionCampaign(appCtx)

	return dClient
}

func (d *etcdDistributedClient) put(appCtx context.Context, key string, value string) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	ctx, cancel := context.WithTimeout(appCtx, EtcdPutGetTimeout)
	defer cancel()

	_, err := d.client.Put(ctx, key, value, clientv3.WithLease(d.session.Lease()))

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

func (d *etcdDistributedClient) runPublishTask(appCtx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	defer close(d.publishChan)

	for {
		select {
		case <-appCtx.Done():
			logger.Infow("Stopping follower task goroutine")
			return

		case message := <-d.publishChan:
			d.put(appCtx, message.key, message.value)
		}
	}
}

func (d *etcdDistributedClient) runWatchTask(appCtx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	// TODO: make variables for watched values

	wchWsConnection := d.client.Watch(appCtx, KeyWsConnectionStat, clientv3.WithPrefix())
	wchTotalAccess := d.client.Watch(appCtx, KeyTotalAccessStat, clientv3.WithPrefix())
	wchLeaderName := d.client.Watch(appCtx, KeyLeaderNameStat, clientv3.WithPrefix())

	stop := false
	for !stop {
		select {
		case <-appCtx.Done():
			stop = true
			break

		case watchResponse := <-wchWsConnection:
			if watchResponse.Canceled {
				logger.Errorw("etcd watch channel is about to close", "key", KeyLeaderNameStat)
				stop = true
				break
			}

			if err := watchResponse.Err(); err != nil {
				logger.Errorw("Unknown watch response error", "error", err)
				continue
			}

			// TODO

		case watchResponse := <-wchTotalAccess:
			if watchResponse.Canceled {
				logger.Errorw("etcd watch channel is about to close", "key", KeyTotalAccessStat)
				stop = true
				break
			}

			if err := watchResponse.Err(); err != nil {
				logger.Errorw("Unknown watch response error", "error", err)
				continue
			}

			d.subscribeTotalAccessCount(&watchResponse)

		case watchResponse := <-wchLeaderName:
			if watchResponse.Canceled {
				logger.Errorw("etcd watch channel is about to close", "key", KeyWsConnectionStat)
				stop = true
				break
			}

			if err := watchResponse.Err(); err != nil {
				logger.Errorw("Unknown watch response error", "error", err)
				continue
			}

			d.subscribeLeaderName(&watchResponse)
		}
	}

	logger.Infow("Stopping watch task goroutine")
}

func (d *etcdDistributedClient) subscribeTotalAccessCount(response *clientv3.WatchResponse) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	for _, ev := range response.Events {
		key := fmt.Sprintf("%s", ev.Kv.Key)
		value := fmt.Sprintf("%s", ev.Kv.Value)
		logger.Infow("WatchEvent", "event_type", ev.Type,
			"event_key", key, "event_value", value)

		message, err := websocket.NewTotalAccessCountMessage(value)
		if err != nil {
			logger.Errorw("Failed to build NewTotalAccessCountMessage", "error", err)
			continue
		}

		d.wsManager.Broadcast(message)
	}
}

func (d *etcdDistributedClient) subscribeLeaderName(response *clientv3.WatchResponse) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	for _, ev := range response.Events {
		key := fmt.Sprintf("%s", ev.Kv.Key)
		value := fmt.Sprintf("%s", ev.Kv.Value)
		logger.Infow("WatchEvent", "event_type", ev.Type,
			"event_key", key, "event_value", value)

		message, err := websocket.NewLeaderNameMessage(value)
		if err != nil {
			logger.Errorw("Failed to build NewLeaderNameMessage", "error", err)
			continue
		}

		d.wsManager.Broadcast(message)
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

		if d.leader == leader {
			return
		}

		d.leader = leader
		d.put(appCtx, KeyLeaderNameStat, d.leader)

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
	d.leader = d.serverName
	d.put(appCtx, KeyLeaderNameStat, d.leader)
}

func (d *etcdDistributedClient) Publish(message *DistributedMessage) {
	d.publishChan <- message
}

func (d *etcdDistributedClient) Stop() {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

	d.session.Close()
	d.client.Close()

	logger.Info("Closed etcd client connection")
}
