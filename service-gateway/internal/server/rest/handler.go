package rest

import (
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations"
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations/access"
	restmodel "github.com/1ambda/go-ref/service-gateway/pkg/api/model"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
	"github.com/jinzhu/gorm"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/model"
	"github.com/go-openapi/swag"
	"time"
)

func ConfigureAPI(db *gorm.DB, api *operations.GatewayAPI) {
	api.AccessAddOneHandler = buildAccessOneHandler(db)
	api.AccessFindOneHandler = buildAccessFindOneHandler(db)
	api.AccessFindAllHandler = buildAccessFindAllHandler(db)
	api.AccessRemoveOneHandler = buildAccessRemoveOneHandler(db)
}

func convertAccessToDbModel(restmodel *restmodel.Access) *model.Access {
	record := model.Access{
		BrowserName:    *restmodel.BrowserName,
		BrowserVersion: *restmodel.BrowserVersion,
		OsName:         *restmodel.OsName,
		OsVersion:      *restmodel.OsVersion,
		IsMobile:       *restmodel.IsMobile,
		Timezone:       *restmodel.Timezone,
		Timestamp:      *restmodel.Timestamp,
		Language:       *restmodel.Language,
		UserAgent:      *restmodel.UserAgent,
	}

	return &record
}

func convertAccessToRestModel(record *model.Access) *restmodel.Access {
	restmodel := restmodel.Access{
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

	return &restmodel
}

func buildAccessOneHandler(db *gorm.DB) access.AddOneHandlerFunc {
	return access.AddOneHandlerFunc(
		func(params access.AddOneParams) middleware.Responder {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			sugar := logger.Sugar()
			sugar.Infow("Creating Access record", "request", params.Body)

			record := convertAccessToDbModel(params.Body)

			if err := db.Create(record).Error; err != nil {
				sugar.Errorw("Failed to create new Access record: %v", "error", err)
				access.NewAddOneDefault(500).WithPayload(&restmodel.Error{
					Code:      500,
					Message:   swag.String(err.Error()),
					Timestamp: time.Now().UTC().String(),
				})
			}

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
				access.NewFindOneDefault(404).WithPayload(&restmodel.Error{
					Code:      404,
					Message:   swag.String(err.Error()),
					Timestamp: time.Now().UTC().String(),
				})
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

			if err := db.Find(&records).Error; err != nil {
				sugar.Errorw("Failed to find all Access records", "error", err)
				access.NewFindAllDefault(500).WithPayload(&restmodel.Error{
					Code:      500,
					Message:   swag.String(err.Error()),
					Timestamp: time.Now().UTC().String(),
				})
			}

			response := make([]*restmodel.Access, 0)
			for _, record := range records {

				response = append(response, &restmodel.Access{
					BrowserName:    &record.BrowserName,
					BrowserVersion: &record.BrowserVersion,
					OsName:         &record.OsName,
					OsVersion:      &record.OsVersion,
					IsMobile:       &record.IsMobile,
					Timezone:       &record.Timezone,
					Timestamp:      &record.Timestamp,
					Language:       &record.Language,
					UserAgent:      &record.UserAgent,
				})
			}

			return access.NewFindAllOK().WithPayload(response)
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
				access.NewAddOneDefault(500).WithPayload(&restmodel.Error{
					Code:      500,
					Message:   swag.String(err.Error()),
					Timestamp: time.Now().UTC().String(),
				})
			}

			return access.NewRemoveOneNoContent()
		})
}
