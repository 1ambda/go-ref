package realtime

import (
	"time"
	"context"

	"go.uber.org/zap"
	"github.com/coreos/etcd/clientv3"
)

type DistributedClient interface {
	attendMembership() error
	renewMembership() error
	run(ctx context.Context)
	Stop()
}

const MembershipTTL = 60 // seconds
const MembershipExtendInterval = 5 * time.Second

type etcdDistributedClient struct {
	client  *clientv3.Client
	leaseId clientv3.LeaseID
}

func NewDistributedClient(appCtx context.Context, endpoints []string) DistributedClient {
	log, _ := zap.NewProduction()
	defer log.Sync()
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

	dClient := &etcdDistributedClient{client: cli}
	if err := dClient.attendMembership(); err != nil {
		logger.Fatalw("Failed to attend membership", "error", err)
	}

	go dClient.run(appCtx)
	go dClient.runRenewMembership(appCtx)

	return dClient
}

func (d *etcdDistributedClient) attendMembership() error {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	response, err := d.client.Grant(context.Background(), MembershipTTL)
	if err != nil {
		return err
	}

	d.leaseId = response.ID

	logger.Infow("Granted etcd client lease", "lease_id", response.ID)

	return nil
}

func (d *etcdDistributedClient) renewMembership() error {
	_, err := d.client.KeepAliveOnce(context.Background(), d.leaseId)
	if err != nil {
		return err
	}

	return nil
}

func (d *etcdDistributedClient) run(ctx context.Context) {
	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	logger.Infow("Running distributed client main loop", "lease_id", d.leaseId)

	for {
		select {
		case <-ctx.Done():
			logger.Infow("Stopping distributed client main loop",
				"lease_id", d.leaseId)
			return

		}
	}
}

func (d *etcdDistributedClient) runRenewMembership(ctx context.Context) {
	extendMembershipTicker := time.NewTicker(MembershipExtendInterval)

	log, _ := zap.NewProduction()
	defer log.Sync()
	logger := log.Sugar()

	for {
		select {
		case <-ctx.Done():
			logger.Infow("Stopping membership renew goroutine for etcd client",
				"lease_id", d.leaseId)
			return

		case <-extendMembershipTicker.C:
			err := d.renewMembership()
			if err != nil {
				logger.Errorw("Failed to keep alive for this etcd client",
					"lease_id", d.leaseId, "error", err)
			}
		}
	}
}

func (d *etcdDistributedClient) Stop() {
	log, _ := zap.NewProduction()
	defer log.Sync() // flushes buffer, if any
	logger := log.Sugar()

	d.client.Close()

	logger.Info("Closed etcd client connection")
}
