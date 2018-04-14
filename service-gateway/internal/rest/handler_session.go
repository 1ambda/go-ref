package rest

import (
	"time"
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/session"
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/internal/config"
	"github.com/1ambda/go-ref/service-gateway/internal/model"
	"errors"
)

const sessionTimeout = 60 * time.Minute

func validateOrGenerateSession(params session.ValidateOrGenerateParams, db *gorm.DB) (*dto.SessionResponse, *dto.Error) {
	sessionId := *params.Body.SessionID

	logger := config.GetLogger()
	logger.Infow("validateOrGenerateSession record", "session", sessionId)


	var session *model.Session = nil
	if sessionId == "" {
		// create new session
		restResp, restErr := createNewSession(db)
		return restResp, restErr

	} else {
		// find existing session
		session = &model.Session{}
		if err := db.Where("session_id = ?", sessionId).First(session).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				// create new session
				restResp, restErr := createNewSession(db)
				return restResp, restErr
			}

			logger.Errorw("Failed to find Session record due to unknown error",
				"session", sessionId, "error", err)
			return nil, buildRestError(err, 500)
		}

		// refresh session if it's expired
		if time.Now().UTC().After(session.ExpiredAt) {
			result := db.Model(&session).Where("session_id = ?", sessionId).Updates(map[string]interface{}{
				"refresh_count": gorm.Expr("refresh_count + ?", 1),
				"refreshed":    true,
				"expired_at":    time.Now().UTC().Add(sessionTimeout),
			})

			if result.Error != nil {
				logger.Errorw("Failed to update Session record due to unknown error",
					"session", sessionId, "error", result.Error)
				restError := buildRestError(result.Error, 500)
				return nil, restError
			}

			if result.RowsAffected < 1 {
				logger.Infow("Failed to find Session record before updating",
					"session", sessionId)
				err := errors.New(gorm.ErrRecordNotFound.Error())
				restError := buildRestError(err, 400)
				return nil, restError
			}
		}
	}

	logger.Infow("Session", "refreshed", session.Refreshed,
		"refreshed_count", session.RefreshCount)

	return model.ConvertToSessionDTO(session), nil
}

func createNewSession(db *gorm.DB) (*dto.SessionResponse, *dto.Error) {
	logger := config.GetLogger()

	session := &model.Session{
		SessionID:    uuid.NewV4().String(),
		ExpiredAt:    time.Now().UTC().Add(sessionTimeout),
		RefreshCount: 0,
		Refreshed:    false,
	}

	if err := db.Create(session).Error; err != nil {
		logger.Errorw("Failed to create Session record: %v", "error", err)
		restError := buildRestError(err, 500)
		return nil, restError
	}

	return model.ConvertToSessionDTO(session), nil
}
