package rest

import (
	"github.com/1ambda/go-ref/service-gateway/internal/config"
	"github.com/1ambda/go-ref/service-gateway/internal/distributed"
	"github.com/1ambda/go-ref/service-gateway/internal/model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/access"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func createAccessRecord() *model.BrowserHistory {
	record := model.BrowserHistory{
		BrowserName:    "Chrome",
		BrowserVersion: "67",
		OsName:         "OSX",
		OsVersion:      "11",
		IsMobile:       "N",
		Timezone:       "KST",
		Timestamp:      "1521919066",
		Language:       "KR",
		UserAgent:      "agent",
	}

	return &record
}

func createAccessRequest() *rest_model.Access {
	record := createAccessRecord()
	return convertAccessToRestModel(record)
}

func insertAccessRecord(db *gorm.DB) *model.BrowserHistory {
	record := createAccessRecord()

	err := db.Create(record).Error
	Expect(err).NotTo(HaveOccurred())

	return record
}

func convertToAccessRequest(record *model.BrowserHistory) *rest_model.Access {
	return &rest_model.Access{
		UUID:           record.UUID,
		BrowserName:    &record.BrowserName,
		BrowserVersion: &record.BrowserVersion,
		OsName:         &record.OsName,
		OsVersion:      &record.OsVersion,
		IsMobile:       &record.IsMobile,
		Timezone:       &record.Timezone,
		Timestamp:      &record.Timestamp,
		Language:       &record.Language,
		UserAgent:      &record.UserAgent,
	}
}

var _ = Describe("Rest: BrowserHistory", func() {
	spec := config.Spec
	spec.Env = config.ENV_TEST

	var (
		db                    *gorm.DB
		ctrl                  *gomock.Controller
		mockDistributedClient *distributed.MockDistributedClient
	)

	BeforeEach(func() {
		db = model.GetDatabase(spec)
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
		db.Close()
	})

	Describe("addOneAccess", func() {
		Context("When got valid params", func() {
			It("should create record", func() {
				params := access.AddOneParams{
					Body: createAccessRequest(),
				}

				mockDistributedClient = distributed.NewMockDistributedClient(ctrl)
				expectedDMessage := distributed.NewBrowserHistoryCountMessage("1")
				mockDistributedClient.EXPECT().Publish(expectedDMessage).Times(1)

				restResp, restErr := addOneAccess(params, db, mockDistributedClient)

				Expect(restErr).To(BeNil())
				Expect(restResp).NotTo(BeNil())
				Expect(restResp.UUID).NotTo(Equal(""))
			})
		})
	})

	Describe("findAllAccess", func() {
		Context("When got valid params", func() {
			It("should return rows", func() {
				insertAccessRecord(db)
				insertAccessRecord(db)
				insertAccessRecord(db)

				var currentPageOffset int32 = 0
				var itemCountPerPage int64 = 2
				params := access.FindAllParams{
					CurrentPageOffset: &currentPageOffset,
					ItemCountPerPage:  &itemCountPerPage,
				}

				pagination, rows, restErr := findAllAccess(params, db)

				Expect(pagination).NotTo(BeNil())
				Expect(len(rows)).To(Equal(2))
				Expect(restErr).To(BeNil())
			})
		})
	})

	Describe("findOneAccess", func() {
		Context("When got invalid ID", func() {
			It("should return rest error", func() {
				params := access.FindOneParams{ID: 0}

				restResp, restErr := findOneAccess(params, db)

				Expect(restResp).To(BeNil())
				Expect(restErr.Code).To(Equal(int64(404)))
			})
		})

		Context("When got valid ID", func() {
			It("should find BrowserHistory record", func() {
				record := insertAccessRecord(db)

				params := access.FindOneParams{ID: int64(record.Id)}
				restResp, restErr := findOneAccess(params, db)

				Expect(restErr).To(BeNil())
				Expect(restResp.ID).To(Equal(int64(record.Id)))
			})
		})
	})

	Describe("removeOneAccess", func() {
		Context("When got invalid ID", func() {
			It("should return rest error", func() {
				params := access.RemoveOneParams{ID: -1}

				restErr := removeOneAccess(params, db)

				Expect(restErr).NotTo(BeNil())
				Expect(restErr.Code).To(Equal(int64(404)))
			})
		})

		Context("When got valid ID", func() {
			It("should delete BrowserHistory record", func() {
				record := insertAccessRecord(db)

				params := access.RemoveOneParams{ID: int64(record.Id)}
				restErr := removeOneAccess(params, db)

				Expect(restErr).To(BeNil())
			})
		})
	})

	Describe("updateOneHandler", func() {
		Context("When got invalid ID", func() {
			It("should return rest error", func() {
				request := convertToAccessRequest(createAccessRecord())
				params := access.UpdateOneParams{ID: -1, Body: request}

				restResp, restErr := updateOneAccess(params, db)

				Expect(restResp).To(BeNil())
				Expect(restErr).NotTo(BeNil())
				Expect(restErr.Code).To(Equal(int64(404)))
			})
		})

		Context("When got valid ID", func() {
			It("should update BrowserHistory record", func() {
				record := insertAccessRecord(db)
				request := convertToAccessRequest(record)

				// update agent value
				newAgent := "agent2"
				request.UserAgent = &newAgent

				params := access.UpdateOneParams{ID: int64(record.Id), Body: request}
				restResp, restErr := updateOneAccess(params, db)

				Expect(restErr).To(BeNil())
				Expect(restResp).NotTo(BeNil())
				Expect(*restResp.UserAgent).To(Equal(newAgent))
			})
		})
	})

})
