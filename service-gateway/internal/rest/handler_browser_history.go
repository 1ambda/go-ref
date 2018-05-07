package rest

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/1ambda/go-ref/service-gateway/internal/distributed"
	"github.com/1ambda/go-ref/service-gateway/internal/model"
	dto "github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/browser_history"
)

func (ctrl *controllerImpl) addOneBrowserHistory(params browser_history.AddOneParams) (*dto.BrowserHistory, *dto.RestError) {
	db := ctrl.db
	dConnector := ctrl.dConnector
	logger := ctrl.logger

	sessionID, restErr := getSessionCookieForRest(params.HTTPRequest, db)
	if restErr != nil {
		return nil, restErr
	}

	logger.Infow("Creating BrowserHistory record",
		"request", params.Body, "session", sessionID)

	var record model.BrowserHistory
	record.ConvertFromDTO(sessionID, params.Body)

	if err := db.Create(&record).Error; err != nil {
		logger.Errorw("Failed to create new BrowserHistory record", "error", err)
		restError := buildRestError(err, dto.RestErrorTypeInternalServer, 500)
		return nil, restError
	}

	var count int64 = 0
	err := db.Table(model.BrowserHistoryTable).Count(&count).Error
	if err != nil {
		logger.Errorw("Failed to create new BrowserHistory record", "error", err)
		restError := buildRestError(err, dto.RestErrorTypeInternalServer, 500)
		return nil, restError
	}

	stringified := fmt.Sprintf("%d", count)
	dConnector.Publish(distributed.NewBrowserHistoryCountMessage(stringified))

	restResp := record.ConvertToDTO()

	return restResp, nil
}

func (ctrl *controllerImpl) findOneBrowserHistory(params browser_history.FindOneParams) (*dto.BrowserHistory, *dto.RestError) {
	db := ctrl.db
	logger := ctrl.logger

	logger.Infow("Finding BrowserHistory record", "id", params.ID)

	var record model.BrowserHistory

	if err := db.Where("id = ?", params.ID).First(&record).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			logger.Infow("Failed to find BrowserHistory record", "id", params.ID)
			return nil, buildRestError(err, dto.RestErrorTypeRecordDoesNotxist, 400)
		}

		logger.Errorw("Failed to find BrowserHistory record due to unknown error",
			"id", params.ID, "error", err)
		return nil, buildRestError(err, dto.RestErrorTypeInternalServer, 500)
	}

	response := record.ConvertToDTO()
	return response, nil
}

func (ctrl *controllerImpl) findAllBrowserHistory(params browser_history.FindAllParams) (*dto.Pagination, []*dto.BrowserHistory, *dto.RestError) {
	db := ctrl.db
	logger := ctrl.logger
	logger.Info("Finding All BrowserHistory records")

	var records []model.BrowserHistory

	// get pagination params
	var count int64 = 0
	currentPageOffset := params.CurrentPageOffset
	itemCountPerPage := params.ItemCountPerPage
	dbOffset := int(int64(*currentPageOffset) * (*itemCountPerPage))
	dbLimit := int(*itemCountPerPage)

	// get filter params
	filterColumn, filterValue, filterErr := extractFilterParams(params.FilterColumn, params.FilterValue)
	if filterErr != nil {
		restError := buildRestError(filterErr, dto.RestErrorTypeBadFilterRequest, 400)
		return nil, nil, restError
	}

	whereQuery, whereArgs, whereClauseErr := getWhereClauseParams(filterColumn, filterValue)
	if whereClauseErr != nil {
		restError := buildRestError(whereClauseErr, dto.RestErrorTypeBadFilterRequest, 400)
		return nil, nil, restError
	}

	chain := db.Table(model.BrowserHistoryTable)
	if whereQuery != "" && whereArgs != "" {
		chain = chain.Where(whereQuery, whereArgs)
	}

	err := chain.Order("created_at asc").Offset(dbOffset).Limit(dbLimit).Find(&records).Error

	// can't use `db.Count()` w/ Offset() + Limit() + Where() due to gorm bug
	chain = db.Table(model.BrowserHistoryTable)
	if whereQuery != "" && whereArgs != "" {
		chain = chain.Where(whereQuery, whereArgs)
	}
	chain.Count(&count)

	if err != nil {
		logger.Errorw("Failed to find all BrowserHistory records", "error", err)
		restError := buildRestError(err, dto.RestErrorTypeInternalServer, 500)
		return nil, nil, restError
	}

	rows := make([]*dto.BrowserHistory, 0)
	for i := range records {
		record := records[i]
		restmodel := record.ConvertToDTO()
		rows = append(rows, restmodel)
	}

	pagination := dto.Pagination{
		ItemCountPerPage:  itemCountPerPage,
		CurrentPageOffset: currentPageOffset,
		TotalItemCount:    &count,
	}

	return &pagination, rows, nil
}

func (ctrl *controllerImpl) removeOneBrowserHistory(params browser_history.RemoveOneParams) *dto.RestError {
	db := ctrl.db
	logger := ctrl.logger
	logger.Infow("Deleting BrowserHistory record", "id", params.ID)

	// https://github.com/jinzhu/gorm/issues/1380
	// https://github.com/jinzhu/gorm/issues/371
	result := db.Where("id = ?", params.ID).Delete(&model.BrowserHistory{})

	if result.Error != nil {
		logger.Errorw("Failed to delete BrowserHistory record due to unknown error",
			"id", params.ID, "error", result.Error)
		restError := buildRestError(result.Error, dto.RestErrorTypeInternalServer, 500)
		return restError
	}

	if result.RowsAffected < 1 {
		logger.Errorw("Failed to find BrowserHistory record before removing", "id", params.ID)
		err := errors.New(gorm.ErrRecordNotFound.Error())
		restError := buildRestError(err, dto.RestErrorTypeRecordDoesNotxist, 400)
		return restError
	}

	return nil
}

func extractFilterParams(ptrFilterColumn *string, ptrFilterValue *string) (string, string, error) {
	// get filtering params
	filterColumn := ""
	if ptrFilterColumn != nil {
		filterColumn = *ptrFilterColumn
	}
	filterValue := ""
	if ptrFilterValue != nil {
		filterValue = *ptrFilterValue
	}

	trimedColumn := strings.TrimSpace(filterColumn)
	trimedValue := strings.TrimSpace(filterValue)

	if trimedColumn == "" && trimedValue != "" {
		err := errors.New("empty filter column")
		return "", "", err
	}

	return trimedColumn, trimedValue, nil
}

func getWhereClauseParams(filterColumn string, filterValue string) (string, string, error) {
	if filterColumn == "" || filterValue == "" {
		return "", "", nil
	}

	if filterColumn == string(dto.BrowserHistoryFilterTypeRecordID) {
		id, err := strconv.ParseInt(filterValue, 10, 64)
		if err != nil {
			return "", "", err
		}

		query := "id = ?"
		args := fmt.Sprintf("%d", id)
		return query, args, nil
	}

	if filterColumn == string(dto.BrowserHistoryFilterTypeSessionID) {
		query := fmt.Sprintf("%s LIKE ?", "session_id")
		args := fmt.Sprintf("%s%%", filterValue) // starts with

		return query, args, nil
	}

	if filterColumn == string(dto.BrowserHistoryFilterTypeBrowserName) {
		query := fmt.Sprintf("%s LIKE ?", "browser_name")
		args := fmt.Sprintf("%%%s%%", filterValue) // bi-directional `like`

		return query, args, nil
	}

	if filterColumn == string(dto.BrowserHistoryFilterTypeLanguage) {
		query := fmt.Sprintf("%s LIKE ?", "language")
		args := fmt.Sprintf("%%%s%%", filterValue) // bi-directional `like`

		return query, args, nil
	}

	if filterColumn == string(dto.BrowserHistoryFilterTypeClientTimezone) {
		query := fmt.Sprintf("%s LIKE ?", "client_timezone")
		args := fmt.Sprintf("%%%s%%", filterValue) // bi-directional `like`

		return query, args, nil
	}

	if filterColumn == string(dto.BrowserHistoryFilterTypeUserAgent) {
		query := fmt.Sprintf("%s LIKE ?", "user_agent")
		args := fmt.Sprintf("%%%s%%", filterValue) // bi-directional `like`

		return query, args, nil
	}

	return "", "", errors.New("invalid filter type")
}
