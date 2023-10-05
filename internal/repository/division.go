package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentreehole/backend/internal/model"
)

type DivisionRepository interface {
	// ListDivisions 获取所有分区
	ListDivisions(ctx context.Context) (response []*model.Division, err error)

	// GetDivisionByID 获取分区
	GetDivisionByID(ctx context.Context, id int) (response *model.Division, err error)

	// CreateDivision 创建分区
	CreateDivision(ctx context.Context, request *model.Division) (response *model.Division, err error)

	// ModifyDivision 修改分区
	ModifyDivision(ctx context.Context, id int, request *model.Division) (response *model.Division, err error)

	// DeleteDivision 删除分区
	DeleteDivision(ctx context.Context, id int) (err error)
}

type divisionRepository struct {
	Repository
}

func (d *divisionRepository) ListDivisions(ctx context.Context) (response []*model.Division, err error) {

	_, err = d.GetCache(ctx).Get(ctx, "divisions", &response)
	if err == nil {
		return response, nil
	}

	err = d.GetDB(ctx).Find(&response).Error
	if err != nil {
		return nil, err
	}

	err = d.GetCache(ctx).Set(ctx, "divisions", response)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (d *divisionRepository) GetDivisionByID(ctx context.Context, id int) (response *model.Division, err error) {

	err = d.GetDB(ctx).First(&response, id).Error
	if err != nil {
		return nil, err
	}

	return response, err
}

func (d *divisionRepository) CreateDivision(ctx context.Context, request *model.Division) (response *model.Division, err error) {

	err = d.GetDB(ctx).FirstOrCreate(&request, model.Division{Name: request.Name}).Error
	if err != nil {
		return nil, err
	}
	response = request

	var divisions []*model.Division
	_, err = d.GetCache(ctx).Get(ctx, "divisions", &divisions)
	if err != nil {
		if errors.Is(err, fmt.Errorf("entry not found")) {
			return response, err
		}
	}
	divisions = append(divisions, response)

	err = d.GetCache(ctx).Set(ctx, "divisions", divisions)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (d *divisionRepository) ModifyDivision(ctx context.Context, id int, request *model.Division) (response *model.Division, err error) {

	request.ID = id
	err = d.GetDB(ctx).Model(&request).Updates(request).Error
	if err != nil {
		return nil, err
	}

	response = request

	return response, err
}

func (d *divisionRepository) DeleteDivision(ctx context.Context, id int) (err error) {

	return d.GetDB(ctx).Delete(&model.Division{}, id).Error
}

func NewDivisionRepository(repository Repository) DivisionRepository {
	return &divisionRepository{Repository: repository}
}
