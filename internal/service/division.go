package service

import (
	"context"
	"github.com/opentreehole/backend/internal/model"
	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
)

type DivisionService interface {
	Service
	// ListDivisions 获取所有分区
	ListDivisions(ctx context.Context) (response []*schema.DivisionResponse, err error)

	// GetDivision 获取分区
	GetDivision(ctx context.Context, id int) (response *schema.DivisionResponse, err error)

	// CreateDivision 创建分区
	CreateDivision(ctx context.Context, request *schema.DivisionCreateRequest) (response *schema.DivisionResponse, err error)

	// ModifyDivision 修改分区
	ModifyDivision(ctx context.Context, id int, request *schema.DivisionModifyRequest) (response *schema.DivisionResponse, err error)

	// DeleteDivision 删除分区
	DeleteDivision(ctx context.Context, id int, request *schema.DivisionDeleteRequest) (response *schema.DivisionResponse, err error)
}

type divisionService struct {
	Service
	Repository repository.DivisionRepository
}

func (d *divisionService) ListDivisions(ctx context.Context) (response []*schema.DivisionResponse, err error) {

	divisionsModel, err := d.Repository.ListDivisions(ctx)
	if err != nil {
		return nil, err
	}

	response = make([]*schema.DivisionResponse, len(divisionsModel))
	for i, model := range divisionsModel {
		response[i] = new(schema.DivisionResponse).FromModel(model, nil)
	}

	return response, nil
}

func (d *divisionService) GetDivision(ctx context.Context, id int) (response *schema.DivisionResponse, err error) {

	divisionModel, err := d.Repository.GetDivisionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response = new(schema.DivisionResponse).FromModel(divisionModel, nil)

	return response, nil
}

func (d *divisionService) CreateDivision(ctx context.Context, request *schema.DivisionCreateRequest) (response *schema.DivisionResponse, err error) {

	divisionModel, err := d.Repository.CreateDivision(ctx, &model.Division{
		Name:        request.Name,
		Description: request.Description,
	})

	response = new(schema.DivisionResponse).FromModel(divisionModel, nil)
	return response, nil
}

func (d *divisionService) ModifyDivision(ctx context.Context, id int, request *schema.DivisionModifyRequest) (response *schema.DivisionResponse, err error) {

	divisionModel, err := d.Repository.ModifyDivision(ctx, id, &model.Division{
		Name:        request.Name,
		Description: request.Description,
		Pinned:      request.Pinned,
	})

	// TODO get Pinned holes

	response = new(schema.DivisionResponse).FromModel(divisionModel, []struct{}{})
	return response, nil
}

func (d *divisionService) DeleteDivision(ctx context.Context, id int, request *schema.DivisionDeleteRequest) (response *schema.DivisionResponse, err error) {

	// TODO move all holes to the target division

	err = d.Repository.DeleteDivision(ctx, id)
	if err != nil {
		return nil, err
	}

	response = &schema.DivisionResponse{}
	return response, nil
}

func NewDivisionService(service Service, repository repository.DivisionRepository) DivisionService {
	return &divisionService{Service: service, Repository: repository}
}
