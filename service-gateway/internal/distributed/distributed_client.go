package distributed

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
	"github.com/1ambda/go-ref/service-gateway/internal/websocket"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"google.golang.org/grpc"
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

	wsManager websocket.WebSocketManager
}

func NewDistributedClient(appCtx context.Context, endpoints []string,
	serverName string, wsManager websocket.WebSocketManager) DistributedClient {
	logger := config.GetLogger()

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
		client: etcdClient, leader: "", serverName: serverName, wsManager: wsManager,
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

func (d *etcdDistributedClient) Publish(message *DistributedMessage) {
	d.publishChan <- message
}

func (d *etcdDistributedClient) Stop() {
	logger := config.GetLogger()

	d.session.Close()
	d.client.Close()

	logger.Info("Closed etcd client connection")
}

func (d *etcdDistributedClient) runPublishTask(appCtx context.Context) {
	logger := config.GetLogger()

	defer close(d.publishChan)

	wsConnCountChan := d.wsManager.SubscribeConnectionCount()
	wsConnCountKey := fmt.Sprintf("%s%s", RangeKeyPrefixWebSocket, config.Spec.ServerName)

	for {
		select {
		case <-appCtx.Done():
			logger.Infow("Stopping follower task goroutine")
			return

		case message := <-d.publishChan:
			d.put(appCtx, message.key, message.value)

		case wsConnCount := <-wsConnCountChan:
			d.put(appCtx, wsConnCountKey, wsConnCount)
		}
	}
}

func (d *etcdDistributedClient) runWatchTask(appCtx context.Context) {
	logger := config.GetLogger()

	// TODO: make variables for watched values

	wchWsConnection := d.client.Watch(appCtx, RangeKeyPrefixWebSocket, clientv3.WithPrefix())
	wchTotalAccess := d.client.Watch(appCtx, SingleKeyTotalAccessCount, clientv3.WithPrefix())
	wchLeaderName := d.client.Watch(appCtx, SingleKeyLeaderName, clientv3.WithPrefix())

	stop := false
	for !stop {
		select {
		case <-appCtx.Done():
			stop = true
			break

		case watchResponse := <-wchWsConnection:
			if watchResponse.Canceled {
				logger.Errorw("etcd watch channel is about to close", "key", RangeKeyPrefixWebSocket)
				stop = true
				break
			}

			if err := watchResponse.Err(); err != nil {
				logger.Errorw("Unknown watch response error", "error", err)
				continue
			}

			d.subscribeWsConnectionCount(appCtx, &watchResponse)

		case watchResponse := <-wchTotalAccess:
			if watchResponse.Canceled {
				logger.Errorw("etcd watch channel is about to close", "key", SingleKeyTotalAccessCount)
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
				logger.Errorw("etcd watch channel is about to close", "key", SingleKeyWsConnectionCount)
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

func (d *etcdDistributedClient) subscribeWsConnectionCount(appCtx context.Context, response *clientv3.WatchResponse) {
	logger := config.GetLogger()

	ctx, cancel := context.WithTimeout(appCtx, EtcdPutGetTimeout)
	defer cancel()

	// watch for websocket connection count is triggered.
	// retrieve all instances' connection count and sum them
	resp, err := d.client.Get(ctx, RangeKeyPrefixWebSocket, clientv3.WithPrefix())

	if err != nil {
		if err == context.Canceled {
			// grpc balancer calls 'Get' with an inflight client.Close
			logger.Errorw("context is canceled. grpc balancer calls 'Get' with an inflight client.Close", "error", err)
		} else if err == grpc.ErrClientConnClosing {
			// grpc balancer calls 'Get' after client.Close.
			logger.Errorw("grpc balancer calls 'Get' after client.Close.", "error", err)
		} else {
			logger.Errorw("Unknown etcd client Get error", "error", err)
		}

		return
	}

	var wsConnCount int64 = 0
	serverCount := 0
	for _, kv := range resp.Kvs {
		key := fmt.Sprintf("%s", kv.Key)
		value := fmt.Sprintf("%s", kv.Value)

		count, err := strconv.ParseInt(value, 10, 64)
		if err != nil || count < 0 {
			logger.Errorf("Failed to parse websocket connection count", "key", key, "count", count)
			continue
		}

		wsConnCount += count
		serverCount += 1
	}

	message, err := websocket.NewConnectionCountMessage(fmt.Sprintf("%d", wsConnCount))
	d.wsManager.Broadcast(message)
	if err != nil {
		logger.Errorw("Failed to build NewConnectionCountMessage", "error", err)
	}

	message, err = websocket.NewNodeCountMessage(fmt.Sprintf("%d", serverCount))
	d.wsManager.Broadcast(message)
	if err != nil {
		logger.Errorw("Failed to build NewNodeCountMessage", "error", err)
	}
}

func (d *etcdDistributedClient) subscribeTotalAccessCount(response *clientv3.WatchResponse) {
	logger := config.GetLogger()

	for _, ev := range response.Events {
		value := fmt.Sprintf("%s", ev.Kv.Value)
		message, err := websocket.NewTotalAccessCountMessage(value)
		if err != nil {
			logger.Errorw("Failed to build NewTotalAccessCountMessage", "error", err)
			continue
		}

		d.wsManager.Broadcast(message)
	}
}

func (d *etcdDistributedClient) subscribeLeaderName(response *clientv3.WatchResponse) {
	logger := config.GetLogger()

	for _, ev := range response.Events {
		value := fmt.Sprintf("%s", ev.Kv.Value)
		message, err := websocket.NewLeaderNameMessage(value)
		if err != nil {
			logger.Errorw("Failed to build NewLeaderNameMessage", "error", err)
			continue
		}

		d.wsManager.Broadcast(message)
	}
}

func (d *etcdDistributedClient) runElectionCampaign(appCtx context.Context) {
	logger := config.GetLogger()

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
	logger := config.GetLogger()

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
		d.put(appCtx, SingleKeyLeaderName, d.leader)

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
	d.put(appCtx, SingleKeyLeaderName, d.leader)
}

func (d *etcdDistributedClient) put(appCtx context.Context, key string, value string) {
	logger := config.GetLogger()

	ctx, cancel := context.WithTimeout(appCtx, EtcdPutGetTimeout)
	defer cancel()

	_, err := d.client.Put(ctx, key, value, clientv3.WithLease(d.session.Lease()))

	if err != nil {
		switch err {
		case context.Canceled:
			logger.Errorw("ctx is canceled by another routine", "error", err)
		case context.DeadlineExceeded:
			logger.Errorw("ctx is attached with a deadline is exceeded", "error", err)
		case rpctypes.ErrEmptyKey:
			logger.Errorw("Empty etdc key", "error", err)
		default:
			logger.Errorw("bad cluster endpoints, which are not etcd servers", "error", err)
		}
	}
}
