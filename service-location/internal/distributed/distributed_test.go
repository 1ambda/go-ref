package distributed_test

import (
	"context"
	"fmt"

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
				cfg := &integration.ClusterConfig{Size: 1, SkipCreatingClient: true}
				cluster := integration.NewClusterV3(T, cfg)
				defer cluster.Terminate(T)

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				connector, err := distributed.New(ctx, cluster.URLs(), "0")
				Expect(connector).NotTo(BeNil())
				Expect(err).To(BeNil())
			})
		})

		Context("without running etcd cluster", func() {
			It("should return error after retry", func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				connector, err := distributed.New(ctx, []string{"http://127.9.9.1:23790"}, "0")
				Expect(connector).To(BeNil())
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("Connector.Stop()", func() {
		It("should stop goroutines and return", func() {
			// given: cluster
			cfg := &integration.ClusterConfig{Size: 1, SkipCreatingClient: true}
			cluster := integration.NewClusterV3(T, cfg)
			defer cluster.Terminate(T)

			// given: connector
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			connector, _ := distributed.New(ctx, cluster.URLs(), "0")

			// when + then
			connector.Stop()
		}, 1)
	})

	Describe("Connector.Publish()", func() {
		It("should send 'Message' to etcd cluster", func() {
			// given: cluster
			cfg := &integration.ClusterConfig{Size: 1, SkipCreatingClient: false}
			cluster := integration.NewClusterV3(T, cfg)
			defer cluster.Terminate(T)

			// given: connector
			ctx, _ := context.WithCancel(context.Background())
			//ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			//defer cancel()
			endpotins := []string{cluster.Members[0].GRPCAddr()}
			fmt.Println(endpotins)
			connector, _ := distributed.New(ctx, endpotins, "0")

			// when
			msg := distributed.Message{Key: "K", Value: "V"}
			err := connector.Publish(ctx, &msg)

			// then
			Expect(err).To(BeNil())

			client := cluster.RandClient()
			resp, err := client.Get(ctx, msg.Key)
			Expect(len(resp.Kvs)).To(Equal(1))
			Expect(string(resp.Kvs[0].Value)).To(Equal(msg.Value))
		})
	})

})
