package rest

import (
	"fmt"
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_model"
	"github.com/1ambda/go-ref/service-gateway/pkg/generated/swagger/rest_server/rest_api/access"
	"github.com/1ambda/go-ref/service-gateway/internal/distributed"
	"github.com/1ambda/go-ref/service-gateway/internal/config"
	"github.com/1ambda/go-ref/service-gateway/internal/model"
)

func addOneAccess(params access.AddOneParams, db *gorm.DB, dClient distributed.DistributedClient) (*rest_model.Access, *rest_model.Error) {
	logger := config.GetLogger()
	logger.Infow("Creating Access record", "request", params.Body)

	record := model.ConvertFromAccessDTO(params.Body)

	if err := db.Create(record).Error; err != nil {
		logger.Errorw("Failed to create new Access record: %v", "error", err)
		restError := buildRestError(err, 500)
		return nil, restError
	}

	var count int64 = 0
	err := db.Table(model.AccessTable).Count(&count).Error
	if err != nil {
		logger.Errorw("Failed to create new Access record: %v", "error", err)
		restError := buildRestError(err, 500)
		return nil, restError
	}

	stringified := fmt.Sprintf("%d", count)
	dClient.Publish(distributed.NewTotalAccessCountMessage(stringified))

	restResp := model.ConvertToAccessDTO(record)

	return restResp, nil
}

func findOneAccess(params access.FindOneParams, db *gorm.DB) (*rest_model.Access, *rest_model.Error) {
	logger := config.GetLogger()
	logger.Infow("Finding Access record", "id", params.ID)

	var record model.Access

	if err := db.Where("id = ?", params.ID).First(&record).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			logger.Infow("Failed to find Access record", "id", params.ID)
			return nil, buildRestError(err, 404)
		}

		logger.Errorw("Failed to find Access record due to unknown error",
			"id", params.ID, "error", err)
		return nil, buildRestError(err, 500)
	}

	response := model.ConvertToAccessDTO(&record)
	return response, nil
}

func findAllAccess(params access.FindAllParams, db *gorm.DB) (*rest_model.Pagination, []*rest_model.Access, *rest_model.Error) {
	logger := config.GetLogger()
	logger.Info("Finding All Access records")

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
		logger.Errorw("Failed to find all Access records", "error", err)
		restError := buildRestError(err, 500)
		return nil, nil, restError
	}

	rows := make([]*rest_model.Access, 0)
	for i := range records {
		record := records[i]
		restmodel := model.ConvertToAccessDTO(&record)
		rows = append(rows, restmodel)
	}

	pagination := rest_model.Pagination{
		ItemCountPerPage:  itemCountPerPage,
		CurrentPageOffset: currentPageOffset,
		TotalItemCount:    &count,
	}

	return &pagination, rows, nil
}

func removeOneAccess(params access.RemoveOneParams, db *gorm.DB) *rest_model.Error {
	logger := config.GetLogger()
	logger.Infow("Deleting Access record", "id", params.ID)

	// https://github.com/jinzhu/gorm/issues/1380
	// https://github.com/jinzhu/gorm/issues/371
	result := db.Where("id = ?", params.ID).Delete(&model.Access{})

	if result.Error != nil {
		logger.Errorw("Failed to delete Access record due to unknown error",
			"id", params.ID, "error", result.Error)
		restError := buildRestError(result.Error, 500)
		return restError
	}

	if result.RowsAffected < 1 {
		logger.Infow("Failed to find Access record before removing", "id", params.ID)
		err := errors.New(gorm.ErrRecordNotFound.Error())
		restError := buildRestError(err, 404)
		return restError
	}

	return nil
}

func updateOneAccess(params access.UpdateOneParams, db *gorm.DB) (*rest_model.Access, *rest_model.Error) {
	logger := config.GetLogger()
	logger.Infow("Updating Access record", "id", params.ID)

	record := model.ConvertFromAccessDTO(params.Body)
	var updated model.Access

	// https://github.com/jinzhu/gorm/issues/891
	result := db.Model(&updated).Where("id = ?", params.ID).Update(record)

	if result.Error != nil {
		logger.Errorw("Failed to update Access record due to unknown error",
			"id", params.ID, "error", result.Error)
		restError := buildRestError(result.Error, 500)
		return nil, restError
	}

	if result.RowsAffected < 1 {
		logger.Infow("Failed to find Access record before updating", "id", params.ID)
		err := errors.New(gorm.ErrRecordNotFound.Error())
		restError := buildRestError(err, 404)
		return nil, restError
	}

	response := model.ConvertToAccessDTO(&updated)
	return response, nil
}