package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

	"github.com/1ambda/go-ref/service-gateway/internal/config"
	"github.com/1ambda/go-ref/service-gateway/internal/model"
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/session"
)


func validateOrGenerateSession(params session.ValidateOrGenerateParams, db *gorm.DB) (*dto.SessionResponse, *dto.Error) {
	sessionId := *params.Body.SessionID

	logger := config.GetLogger()
	logger.Infow("validateOrGenerateSession record", "session", sessionId)

	var record *model.Session
	var restErr *dto.Error
	if sessionId == "" {
		// create new session
		record, restErr = createNewSession(db)
	} else {
		// find existing session and refresh it if it's expired
		record, restErr = refreshSession(db, sessionId)
	}

	return record.ConvertToSessionDTO(), restErr
}

func createNewSession(db *gorm.DB) (*model.Session, *dto.Error) {
	logger := config.GetLogger()

	record := &model.Session{
		SessionID:    uuid.NewV4().String(),
		ExpiredAt:    time.Now().UTC().Add(config.SessionTimeout),
		RefreshCount: 0,
		Refreshed:    false,
	}

	if err := db.Create(record).Error; err != nil {
		logger.Errorw("Failed to create Session record: %v", "error", err)
		restError := buildRestError(err, 500)
		return nil, restError
	}

	return record, nil
}

func refreshSession(db *gorm.DB, sessionId string) (*model.Session, *dto.Error) {
	logger := config.GetLogger()

	// find existing session
	record := &model.Session{}
	if err := db.Where("session_id = ?", sessionId).First(record).Error; err != nil {
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
	if time.Now().UTC().After(record.ExpiredAt) {
		result := db.Model(&record).Where("session_id = ?", sessionId).Updates(map[string]interface{}{
			"refresh_count": gorm.Expr("refresh_count + ?", 1),
			"refreshed":     true,
			"expired_at":    time.Now().UTC().Add(config.SessionTimeout),
		})

		if result.Error != nil {
			logger.Errorw("Failed to update Session record due to unknown error",
				"session", sessionId, "error", result.Error)
			restError := buildRestError(result.Error, 500)
			return nil, restError
		}

		if result.RowsAffected < 1 {
			logger.Errorw("Failed to find Session record before updating",
				"session", sessionId)
			err := errors.New(gorm.ErrRecordNotFound.Error())
			restError := buildRestError(err, 400)
			return nil, restError
		}
	}

	return record, nil
}

func getSessionCookieForRest(req *http.Request) (string, *dto.Error) {
	cookie, err := req.Cookie(config.SessionKey)

	if err != nil {
		restError := buildRestError(err, 500)
		return "", restError
	}

	if cookie == nil || cookie.Value == "" {
		err := errors.New("empty session cookie")
		restError := buildRestError(err, 400)
		return "", restError
	}

	return cookie.Value, nil
}
