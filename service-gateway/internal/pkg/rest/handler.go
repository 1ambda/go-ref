package rest

import (
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/access"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
	"github.com/jinzhu/gorm"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/model"
	"github.com/go-openapi/swag"
	"time"
	"github.com/satori/go.uuid"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/service"
)

func Configure(db *gorm.DB, api *rest_api.GatewayRestAPI, r *service.RealtimeStatService) {
	api.AccessAddOneHandler = buildAccessAddOneHandler(db, r)
	api.AccessFindOneHandler = buildAccessFindOneHandler(db)
	api.AccessFindAllHandler = buildAccessFindAllHandler(db)
	api.AccessRemoveOneHandler = buildAccessRemoveOneHandler(db)
	api.AccessUpdateOneHandler = buildAccessUpdateOneHandler(db)
}

func convertAccessToDbModel(swaggerModel *rest_model.Access) *model.Access {
	uuid := uuid.NewV4()

	record := model.Access{
		BrowserName:    *swaggerModel.BrowserName,
		BrowserVersion: *swaggerModel.BrowserVersion,
		OsName:         *swaggerModel.OsName,
		OsVersion:      *swaggerModel.OsVersion,
		IsMobile:       *swaggerModel.IsMobile,
		Timezone:       *swaggerModel.Timezone,
		Timestamp:      *swaggerModel.Timestamp,
		Language:       *swaggerModel.Language,
		UserAgent:      *swaggerModel.UserAgent,
		UUID:           uuid.String(),
	}

	return &record
}

func convertAccessToRestModel(record *model.Access) *rest_model.Access {
	swaggerModel := rest_model.Access{
		ID:             int64(record.Id),
		BrowserName:    &record.BrowserName,
		BrowserVersion: &record.BrowserVersion,
		OsName:         &record.OsName,
		OsVersion:      &record.OsVersion,
		IsMobile:       &record.IsMobile,
		Timezone:       &record.Timezone,
		Timestamp:      &record.Timestamp,
		Language:       &record.Language,
		UserAgent:      &record.UserAgent,
		UUID:           record.UUID,
	}

	return &swaggerModel
}

func buildRestError(err error) *rest_model.Error {
	return &rest_model.Error{
		Code:      500,
		Message:   swag.String(err.Error()),
		Timestamp: time.Now().UTC().String(),
	}
}

func buildAccessAddOneHandler(db *gorm.DB, r *service.RealtimeStatService) access.AddOneHandlerFunc {
	return access.AddOneHandlerFunc(
		func(params access.AddOneParams) middleware.Responder {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			sugar := logger.Sugar()
			sugar.Infow("Creating Access record", "request", params.Body)

			record := convertAccessToDbModel(params.Body)

			if err := db.Create(record).Error; err != nil {
				sugar.Errorw("Failed to create new Access record: %v", "error", err)
				restError := buildRestError(err)
				access.NewAddOneDefault(500).WithPayload(restError)
			}

			r.BroadcastToTalAccessCount()

			return access.NewAddOneCreated().WithPayload(params.Body)
		})
}

func buildAccessFindOneHandler(db *gorm.DB) access.FindOneHandlerFunc {
	return access.FindOneHandlerFunc(
		func(params access.FindOneParams) middleware.Responder {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			sugar := logger.Sugar()
			sugar.Infow("Finding Access record", "id", params.ID)

			var record model.Access

			if err := db.Where("id = ?", params.ID).First(&record).Error; err != nil {
				sugar.Errorw("Failed to create new Access record", "error", err)
				restError := buildRestError(err)
				access.NewFindOneDefault(404).WithPayload(restError)
			}

			response := convertAccessToRestModel(&record)
			return access.NewFindOneOK().WithPayload(response)
		})
}

func buildAccessFindAllHandler(db *gorm.DB) access.FindAllHandlerFunc {
	return access.FindAllHandlerFunc(
		func(params access.FindAllParams) middleware.Responder {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			sugar := logger.Sugar()
			sugar.Info("Finding All Access records")

			var records []model.Access
			var count int64 = 0
			currentPageOffset := params.CurrentPageOffset
			itemCountPerPage := params.ItemCountPerPage

			dbOffset := int64(*currentPageOffset) * (*itemCountPerPage)

			err := db.
				Table(model.AccessTable).
				Count(&count).
				Offset(int(dbOffset)).
				Limit(int(*itemCountPerPage)).
				Find(&records).
				Error
			if err != nil {
				sugar.Errorw("Failed to find all Access records", "error", err)
				restError := buildRestError(err)
				access.NewFindAllDefault(500).WithPayload(restError)
			}

			rows := make([]*rest_model.Access, 0)
			for i, _ := range records {
				record := records[i]
				restmodel := convertAccessToRestModel(&record)
				rows = append(rows, restmodel)
			}

			pagination := rest_model.Pagination{
				ItemCountPerPage:  itemCountPerPage,
				CurrentPageOffset: currentPageOffset,
				TotalItemCount:    &count,
			}

			return access.NewFindAllOK().WithPayload(&rest_model.FindAllOKBody{
				Pagination: &pagination,
				Rows:       rows,
			})
		})
}

func buildAccessRemoveOneHandler(db *gorm.DB) access.RemoveOneHandlerFunc {
	return access.RemoveOneHandlerFunc(
		func(params access.RemoveOneParams) middleware.Responder {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			sugar := logger.Sugar()
			sugar.Infow("Deleting Access record", "id", params.ID)

			if err := db.Where("id = ?", params.ID).Delete(&model.Access{}).Error; err != nil {
				sugar.Errorw("Failed to delete new Access record: %v", "error", err)
				restError := buildRestError(err)
				access.NewAddOneDefault(500).WithPayload(restError)
			}

			return access.NewRemoveOneNoContent()
		})
}

func buildAccessUpdateOneHandler(db *gorm.DB) access.UpdateOneHandlerFunc {
	return access.UpdateOneHandlerFunc(
		func(params access.UpdateOneParams) middleware.Responder {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			sugar := logger.Sugar()
			sugar.Infow("Updating Access record", "id", params.ID)

			record := convertAccessToDbModel(params.Body)
			var updated model.Access

			if err := db.Model(&updated).Where("id = ?", params.ID).Update(record).Error; err != nil {
				sugar.Errorw("Failed to update new Access record: %v", "error", err)
				restError := buildRestError(err)
				access.NewAddOneDefault(500).WithPayload(restError)
			}

			response := convertAccessToRestModel(&updated)

			return access.NewUpdateOneOK().WithPayload(response)
		})
}
