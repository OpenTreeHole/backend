package repository

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
)

type CourseGroupRepository interface {
	Repository

	FindGroups(ctx context.Context, condition func(db *gorm.DB) *gorm.DB) (groups []*model.CourseGroup, err error)
	FindGroupByID(ctx context.Context, id int, condition func(db *gorm.DB) *gorm.DB) (group *model.CourseGroup, err error)
	FindGroupByCode(ctx context.Context, code string, condition func(db *gorm.DB) *gorm.DB) (group *model.CourseGroup, err error)
	CreateGroup(ctx context.Context, group *model.CourseGroup) (err error)
	FindGroupsWithCourses(ctx context.Context, refresh bool) (groups []*model.CourseGroup, hash string, err error)
}

type courseGroupRepository struct {
	Repository
	courseGroupHash string
}

func NewCourseGroupRepository(repository Repository) CourseGroupRepository {
	return &courseGroupRepository{Repository: repository}
}

/* 接口实现 */

func (r *courseGroupRepository) FindGroups(ctx context.Context, condition func(db *gorm.DB) *gorm.DB) (groups []*model.CourseGroup, err error) {
	groups = make([]*model.CourseGroup, 5)
	err = condition(r.GetDB(ctx)).Find(&groups).Error
	return
}

func (r *courseGroupRepository) FindGroupsWithCourses(ctx context.Context, refresh bool) (groups []*model.CourseGroup, hash string, err error) {
	groups = make([]*model.CourseGroup, 5)
	if !refresh {
		// get from cache
		_, err = r.GetCache(ctx).Get(ctx, "danke:course_group", &groups)
		hash = r.courseGroupHash
	}
	if err != nil || refresh {
		// get from db
		err = r.GetDB(ctx).Preload("Courses").Find(&groups).Error
		if err != nil {
			return nil, "", err
		}

		// set cache
		err = r.GetCache(ctx).Set(ctx, "danke:course_group", groups, store.WithExpiration(24*time.Hour))
		if err != nil {
			return nil, "", err
		}

		// set hash
		data, err := msgpack.Marshal(groups)
		if err != nil {
			return nil, "", err
		}
		hashBytes := sha3.Sum256(data)
		if err != nil {
			return nil, "", err
		}
		hash = base64.RawStdEncoding.EncodeToString(hashBytes[:])
		r.courseGroupHash = hash
	}
	return
}

func (r *courseGroupRepository) FindGroupByID(ctx context.Context, id int, condition func(db *gorm.DB) *gorm.DB) (group *model.CourseGroup, err error) {
	group = new(model.CourseGroup)
	err = condition(r.GetDB(ctx)).First(group, id).Error
	return
}

func (r *courseGroupRepository) FindGroupByCode(ctx context.Context, code string, condition func(db *gorm.DB) *gorm.DB) (group *model.CourseGroup, err error) {
	group = new(model.CourseGroup)
	err = condition(r.GetDB(ctx)).Where("code = ?", code).First(group).Error
	return
}

func (r *courseGroupRepository) CreateGroup(ctx context.Context, group *model.CourseGroup) (err error) {
	return r.GetDB(ctx).Create(group).Error
}
