package location_test

import (
	"context"

	"github.com/1ambda/go-ref/service-location/internal/distributed"
	"github.com/1ambda/go-ref/service-location/internal/location"
	"github.com/1ambda/go-ref/service-location/pkg/generated/grpc"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("location/", func() {

	var (
		ctrl          *gomock.Controller
		mockConnector *distributed.MockConnector
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("LocationServer.Add()", func() {
		Context("When connector returns another leader", func() {
			It("should throw error with InvalidArgument code", func() {
				// given: mock connector
				mockConnector = distributed.NewMockConnector(ctrl)
				mockConnector.EXPECT().GetLeaderOrCampaign(
					gomock.Any(), gomock.Any(),
				).Return("that", nil)

				// given: service
				srvName := "this"
				svc, err := location.New(srvName, mockConnector)
				Expect(err).To(BeNil())

				// when
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				req := &pb.LocationRequest{
					LocationContext: &pb.LocationContext{SessionId: "session", Country: "somewhere",},
				}
				resp, err := svc.Add(ctx, req)

				// then
				Expect(resp).To(BeNil())
				code := status.Code(err)
				Expect(code).To(Equal(codes.InvalidArgument))

			})

		})

		Context("When connector the server name as a leader", func() {
			It("should return response", func() {
				// given: mock connector
				srvName := "this"
				mockConnector = distributed.NewMockConnector(ctrl)
				mockConnector.EXPECT().GetLeaderOrCampaign(
					gomock.Any(), gomock.Any(),
				).Return(srvName, nil)

				// given: service
				svc, err := location.New(srvName, mockConnector)
				Expect(err).To(BeNil())

				// when
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				req := &pb.LocationRequest{
					LocationContext: &pb.LocationContext{SessionId: "session", Country: "somewhere",},
				}
				resp, err := svc.Add(ctx, req)

				// then
				Expect(resp).NotTo(BeNil())
			})
		})

	})
})
