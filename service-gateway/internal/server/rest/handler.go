package rest

import (
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations"
	"github.com/1ambda/go-ref/service-gateway/pkg/api/rest/operations/access"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
	"github.com/jinzhu/gorm"
	"github.com/1ambda/go-ref/service-gateway/internal/pkg/model"
)

func ConfigureAPI(db *gorm.DB, api *operations.GatewayAPI) {
	api.AccessAddOneHandler = access.AddOneHandlerFunc(
		func(params access.AddOneParams) middleware.Responder {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			sugar := logger.Sugar()
			sugar.Infow("Creating Access",
				"access", params.Body,
			)

			record := model.Access{
				BrowserName:    *params.Body.BrowserName,
				BrowserVersion: *params.Body.BrowserVersion,
				OsName:         *params.Body.OsName,
				OsVersion:      *params.Body.OsVersion,
				IsMobile:       *params.Body.IsMobile,
				Timezone:       *params.Body.Timezone,
				Timestamp:      *params.Body.Timestamp,
				Language:       *params.Body.Language,
				UserAgent:      *params.Body.UserAgent,
			}

			db.Create(&record)

			return access.NewAddOneCreated().WithPayload(params.Body)
		})

}
