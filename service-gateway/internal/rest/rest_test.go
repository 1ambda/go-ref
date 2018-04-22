package rest

import (
	"github.com/1ambda/go-ref/service-gateway/internal/config"
	"github.com/1ambda/go-ref/service-gateway/internal/distributed"
	"github.com/1ambda/go-ref/service-gateway/internal/model"
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/browser_history"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func createAccessRecord() *model.BrowserHistory {
	record := model.BrowserHistory{
		BrowserName:     "Chrome",
		BrowserVersion:  "67",
		OsName:          "OSX",
		OsVersion:       "11",
		IsMobile:        false,
		ClientTimezone:  "KST",
		ClientTimestamp: "1521919066",
		Language:        "KR",
		UserAgent:       "agent",
	}

	return &record
}

func createAccessRequest() *dto.BrowserHistory {
	record := createAccessRecord()
	return convertToAccessRequest(record)
}

func insertAccessRecord(db *gorm.DB) *model.BrowserHistory {
	record := createAccessRecord()

	err := db.Create(record).Error
	Expect(err).NotTo(HaveOccurred())

	return record
}

func convertToAccessRequest(record *model.BrowserHistory) *dto.BrowserHistory {
	return &dto.BrowserHistory{
		BrowserName:     &record.BrowserName,
		BrowserVersion:  &record.BrowserVersion,
		OsName:          &record.OsName,
		OsVersion:       &record.OsVersion,
		IsMobile:        &record.IsMobile,
		ClientTimezone:  &record.ClientTimezone,
		ClientTimestamp: &record.ClientTimestamp,
		Language:        &record.Language,
		UserAgent:       &record.UserAgent,
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

	Describe("addOne", func() {
		Context("When got valid params", func() {
			It("should create record", func() {
				params := browser_history.AddOneParams{
					Body: createAccessRequest(),
				}

				mockDistributedClient = distributed.NewMockDistributedClient(ctrl)
				expectedDMessage := distributed.NewBrowserHistoryCountMessage("1")
				mockDistributedClient.EXPECT().Publish(expectedDMessage).Times(1)

				restResp, restErr := addOneBrowserHistory(params, db, mockDistributedClient)

				Expect(restErr).To(BeNil())
				Expect(restResp).NotTo(BeNil())
				Expect(restResp.RecordID).To(Equal(int64(1)))
			})
		})
	})

	Describe("findAll", func() {
		Context("When got valid params", func() {
			It("should return rows", func() {
				insertAccessRecord(db)
				insertAccessRecord(db)
				insertAccessRecord(db)

				var currentPageOffset int32 = 0
				var itemCountPerPage int64 = 2
				params := browser_history.FindAllParams{
					CurrentPageOffset: &currentPageOffset,
					ItemCountPerPage:  &itemCountPerPage,
				}

				pagination, rows, restErr := findAllBrowserHistory(params, db)

				Expect(pagination).NotTo(BeNil())
				Expect(len(rows)).To(Equal(2))
				Expect(restErr).To(BeNil())
			})
		})
	})

	Describe("findOne", func() {
		Context("When got invalid ID", func() {
			It("should return rest error", func() {
				params := browser_history.FindOneParams{ID: 0}

				restResp, restErr := findOneBrowserHistory(params, db)

				Expect(restResp).To(BeNil())
				Expect(restErr.Code).To(Equal(int64(404)))
			})
		})

		Context("When got valid ID", func() {
			It("should find record", func() {
				record := insertAccessRecord(db)

				params := browser_history.FindOneParams{ID: int64(record.ID)}
				restResp, restErr := findOneBrowserHistory(params, db)

				Expect(restErr).To(BeNil())
				Expect(restResp.RecordID).To(Equal(int64(record.ID)))
			})
		})
	})

	Describe("removeOne", func() {
		Context("When got invalid ID", func() {
			It("should return rest error", func() {
				params := browser_history.RemoveOneParams{ID: -1}

				restErr := removeOneBrowserHistory(params, db)

				Expect(restErr).NotTo(BeNil())
				Expect(restErr.Code).To(Equal(int64(404)))
			})
		})

		Context("When got valid ID", func() {
			It("should delete record", func() {
				record := insertAccessRecord(db)

				params := browser_history.RemoveOneParams{ID: int64(record.ID)}
				restErr := removeOneBrowserHistory(params, db)

				Expect(restErr).To(BeNil())
			})
		})
	})
})
