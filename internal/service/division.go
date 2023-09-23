package service

import (
	"context"

	"github.com/opentreehole/backend/internal/schema"
)

type DivisionService interface {
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
}

func (d *divisionService) ListDivisions(ctx context.Context) (response []*schema.DivisionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *divisionService) GetDivision(ctx context.Context, id int) (response *schema.DivisionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *divisionService) CreateDivision(ctx context.Context, request *schema.DivisionCreateRequest) (response *schema.DivisionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *divisionService) ModifyDivision(ctx context.Context, id int, request *schema.DivisionModifyRequest) (response *schema.DivisionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *divisionService) DeleteDivision(ctx context.Context, id int, request *schema.DivisionDeleteRequest) (response *schema.DivisionResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func NewDivisionService(service Service) DivisionService {
	return &divisionService{Service: service}
}
