package repository

import (
	"context"

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
	//TODO implement me
	panic("implement me")
}

func (d *divisionRepository) GetDivisionByID(ctx context.Context, id int) (response *model.Division, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *divisionRepository) CreateDivision(ctx context.Context, request *model.Division) (response *model.Division, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *divisionRepository) ModifyDivision(ctx context.Context, id int, request *model.Division) (response *model.Division, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *divisionRepository) DeleteDivision(ctx context.Context, id int) (err error) {
	//TODO implement me
	panic("implement me")
}

func NewDivisionRepository(repository Repository) DivisionRepository {
	return &divisionRepository{Repository: repository}
}
