package model

import (
	"context"
	"encoding/base64"
	"regexp"
	"time"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/opentreehole/backend/common"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"
)

// CourseGroup 课程组
type CourseGroup struct {
	ID          int        `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Name        string     `json:"name" gorm:"not null"`                   // 课程组名称
	Code        string     `json:"code" gorm:"not null;index:,length:6"`   // 课程组编号
	Credits     []float64  `json:"credits" gorm:"serializer:json"`         // 学分
	Department  string     `json:"department" gorm:"not null"`             // 开课学院
	CampusName  string     `json:"campus_name" gorm:"not null"`            // 开课校区
	CourseCount int        `json:"course_count" gorm:"not null;default:0"` // 课程数量
	ReviewCount int        `json:"review_count" gorm:"not null;default:0"` // 评价数量
	Courses     CourseList `json:"courses"`
	Teachers []*Teacher    `gorm:"many2many:teacher_course_groups;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

var CourseCodeRegexp = regexp.MustCompile(`^([A-Z]{3,})([0-9]{2,})`)

var courseGroupHash string

func FindGroupsWithCourses(refresh bool) (groups []*CourseGroup, hash string, err error) {
	groups = make([]*CourseGroup, 5)
	if !refresh {
		// get from cache
		_, err = common.Cache.Get(context.Background(), "danke:course_group", &groups)
		hash = courseGroupHash
	}
	if err != nil || refresh {
		// get from db
		err = DB.Preload("Courses").Find(&groups).Error
		if err != nil {
			return nil, "", err
		}

		// set cache
		err = common.Cache.Set(context.Background(), "danke:course_group", groups, store.WithExpiration(24*time.Hour))
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
		courseGroupHash = hash
	}
	return
}

type CourseGroupList []*CourseGroup
