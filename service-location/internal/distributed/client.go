package distributed

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/1ambda/go-ref/service-location/internal/config"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/pkg/errors"
)

const CampaignInterval = 5 * time.Second
const ElectionTimeout = 10 * time.Second
const EtcdSessionTTL = 120 // second
const EtcdPutGetTimeout = 5 * time.Second
const EtcdDialTimeout = 3 * time.Second

const ElectionPathPrefix = "service-location"

// distributed storage connector.
type Connector interface {
	GetLeaderOrCampaign(electSubPath string, electProclaim string) (string, error)
	Publish(ctx context.Context, message *Message) error
	Stop()
}

const messageChanBufferSize = 100

type etcdConnector struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	lock   sync.RWMutex

	etcdClient *clientv3.Client
	serverName string
	endpoints  []string

	messageChan chan *Message
}

type LeaderChangeHandler func(string)

func New(appCtx context.Context, endpoints []string, serverName string) (Connector, error) {
	logger := config.GetLogger()

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: EtcdDialTimeout,
	})

	if err != nil {
		if err == context.DeadlineExceeded {
			err = errors.Wrap(err, "Failed to connect to etcd cluster (DialTimeout)")
			return nil, err
		}

		err = errors.Wrap(err, "Failed to connect to etcd cluster (Unknown)")
		return nil, err
	}

	ctx, cancel := context.WithCancel(appCtx)

	connector := &etcdConnector{
		ctx:    ctx,
		cancel: cancel,

		etcdClient: etcdClient,
		endpoints:  endpoints,
		serverName: serverName,

		messageChan: make(chan *Message, messageChanBufferSize),
	}

	connector.wg.Add(1)
	go connector.run()

	logger.Info("Starting Connector")

	return connector, nil
}

func (c *etcdConnector) Publish(ctx context.Context, m *Message) error {
	putCtx, cancel := context.WithTimeout(c.ctx, EtcdPutGetTimeout)
	defer cancel()

	key := ""

	if m.SubKey == "" {
		key = fmt.Sprintf("%s/%s/value", ElectionPathPrefix, m.Key)
	} else {
		key = fmt.Sprintf("%s/%s/value/%s", ElectionPathPrefix, m.Key, m.SubKey)
	}

	_, err := c.etcdClient.Put(putCtx, key, m.Value)

	if err != nil {
		switch err {
		case context.Canceled:
			err = errors.Wrap(err, "Failed to put etcd KV (Canceled)")
		case context.DeadlineExceeded:
			err = errors.Wrap(err, "Failed to put etcd KV (DeadlineExceeded)")
		case rpctypes.ErrEmptyKey:
			err = errors.Wrap(err, "Failed to put etcd KV (ErrEmptyKey)")
		default:
			err = errors.Wrap(err, "Failed to put etcd KV (Unknown)")
		}
	}

	return err
}

func (c *etcdConnector) Stop() {
	c.cancel()
	c.wg.Wait()
	close(c.messageChan)
}

func (c *etcdConnector) run() {
	logger := config.GetLogger()

	for {
		select {
		case <-c.ctx.Done():
			logger.Infow("Stopping etcd connector")
			c.wg.Done()
			return

			//case message := <- c.messageChan:
			//	logger.Info("Sending message")
			//	err := c.put(message)
			//	if err != nil {
			//		logger.Errorw("Failed to send Message", "err", err)
			//	}

		}
	}
}

func (c *etcdConnector) GetLeaderOrCampaign(electSubPath string, electProclaim string) (string, error) {
	session, err := concurrency.NewSession(c.etcdClient)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get etcd session")
	}

	electPath := fmt.Sprintf("%s/%s/leader", ElectionPathPrefix, electSubPath)
	elect := concurrency.NewElection(session, electPath)
	ctx, cancel := context.WithTimeout(c.ctx, ElectionTimeout)
	defer cancel()

	resp, err := elect.Leader(ctx)

	if err == nil {
		kv := resp.Kvs[0]
		leader := string(kv.Value)
		return leader, nil
	}

	if err != nil && err != concurrency.ErrElectionNoLeader {
		message := fmt.Sprintf("Failed to get leader for path: %s", electPath)
		return "", errors.Wrap(err, message)
	}

	// err != nil and err == concurrency.ErrElectionNoLeader, will campaign
	if err := elect.Campaign(ctx, electProclaim); err != nil {
		message := fmt.Sprintf("Failed to campign leader for path: %s value: %s",
			electPath, electProclaim)
		return "", errors.Wrap(err, message)
	}

	return electProclaim, nil
}
