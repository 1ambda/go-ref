package distributed_test

import (
	"context"
	"time"

	"github.com/1ambda/go-ref/service-location/internal/distributed"
	"github.com/coreos/etcd/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Distributed", func() {
	BeforeEach(func() {

	})

	AfterEach(func() {

	})

	Describe("Connector.New()", func() {
		Context("with running etcd cluster", func() {
			It("should return Connector", func() {
				ctx := context.TODO()
				ctx, cancel := context.WithTimeout(ctx, time.Second*30)
				defer cancel()

				cfg := &integration.ClusterConfig{Size: 1, SkipCreatingClient: true}

				cluster := integration.NewClusterV3(T, cfg)
				defer cluster.Terminate(T)

				connector, err := distributed.New(ctx, cluster.URLs(), "0")
				Expect(connector).NotTo(BeNil())
				Expect(err).To(BeNil())

			})
		})

		Context("without running etcd cluster", func() {
			It("should return error after retry", func() {
				ctx := context.TODO()
				ctx, cancel := context.WithTimeout(ctx, time.Second*5)
				defer cancel()

				connector, err := distributed.New(ctx, []string{"http://127.9.9.1:23790"}, "0")
				Expect(connector).To(BeNil())
				Expect(err).NotTo(BeNil())
			})
		})
	})


})


