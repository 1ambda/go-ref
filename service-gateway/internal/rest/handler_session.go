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
	"net/http"
)

const sessionTimeout = 60 * time.Minute
const sessionKey = "sessionID"

func validateOrGenerateSession(params session.ValidateOrGenerateParams, db *gorm.DB) (*dto.SessionResponse, *dto.Error) {
	sessionId := *params.Body.SessionID

	logger := config.GetLogger()
	logger.Infow("validateOrGenerateSession record", "session", sessionId)

	var sessionRecord *model.Session
	var restErr *dto.Error
	if sessionId == "" {
		// create new session
		sessionRecord, restErr = createNewSession(db)
	} else {
		// find existing session and refresh it if it's expired
		sessionRecord, restErr = refreshSession(db, sessionId)
	}

	return model.ConvertToSessionDTO(sessionRecord), restErr
}

func createNewSession(db *gorm.DB) (*model.Session, *dto.Error) {
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

	return session, nil
}

func refreshSession(db *gorm.DB, sessionId string) (*model.Session, *dto.Error) {
	logger := config.GetLogger()

	// find existing session
	sessionRecord := &model.Session{}
	if err := db.Where("session_id = ?", sessionId).First(sessionRecord).Error; err != nil {
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
	if time.Now().UTC().After(sessionRecord.ExpiredAt) {
		result := db.Model(&sessionRecord).Where("session_id = ?", sessionId).Updates(map[string]interface{}{
			"refresh_count": gorm.Expr("refresh_count + ?", 1),
			"refreshed":     true,
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

	return sessionRecord, nil
}

func getSessionCookie(req *http.Request) (string, *dto.Error){
	cookie, err := req.Cookie(sessionKey)

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
