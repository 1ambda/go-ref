package distributed

import (
	"context"
	"sync"
	"time"

	"github.com/1ambda/go-ref/service-location/internal/config"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/pkg/errors"
)

const CampaignInterval = 5 * time.Second
const ElectionTimeout = 10 * time.Second
const EtcdSessionTTL = 120 // second
const EtcdPutGetTimeout = 5 * time.Second

const EtcdKeyMessage = "service-location/message"

// distributed storage connector.
type Connector interface {
	Publish(ctx context.Context, message *Message) error
	Stop()
}

const messageChanBufferSize = 100

type etcdConnectorImpl struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	lock sync.RWMutex

	etcdClient *clientv3.Client
	serverName string
	endpoints  []string

	messageChan chan *Message
}

func New(appCtx context.Context, endpoints []string, serverName string) (Connector, error) {
	logger := config.GetLogger()

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

	ctx, cancel := context.WithCancel(appCtx)

	connector := &etcdConnectorImpl{
		ctx:        ctx,
		cancel:     cancel,

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

func (c *etcdConnectorImpl) Publish(ctx context.Context, message *Message) error {

	_, err := c.etcdClient.Put(ctx, message.Key, message.Value)

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

func (c *etcdConnectorImpl) Stop() {
	c.cancel()
	c.wg.Wait()
	close(c.messageChan)
}

func (c *etcdConnectorImpl) run() {
	logger := config.GetLogger()

	for {
		select {
		case <-c.ctx.Done():
			logger.Infow("Stopping etcd connector")
			c.wg.Done()
			return

		case message := <- c.messageChan:
			logger.Info("Sending message")
			err := c.put(message)
			if err != nil {
				logger.Errorw("Failed to send Message", "err", err)
			}

		}
	}
}

func (c *etcdConnectorImpl) put(m *Message) error {
	putCtx, cancel := context.WithTimeout(c.ctx, EtcdPutGetTimeout)
	defer cancel()

	_, err := c.etcdClient.Put(putCtx, m.Key, m.Value)

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

type Message struct {
	Key   string
	Value string
}
