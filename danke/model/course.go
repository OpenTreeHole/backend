package model

import (
	"context"
	"github.com/opentreehole/backend/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Course 课程
type Course struct {
	ID            int          `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Name          string       `json:"name" gorm:"not null"`                  // 课程名称
	Code          string       `json:"code" gorm:"not null"`                  // 课程编号
	CodeID        string       `json:"code_id" gorm:"not null"`               // 选课序号。用于区分同一课程编号的不同平行班
	Credit        float64      `json:"credit" gorm:"not null"`                // 学分
	Department    string       `json:"department" gorm:"not null"`            // 开课学院
	CampusName    string       `json:"campus_name" gorm:"not null"`           // 开课校区
	Teachers      string       `json:"teachers" gorm:"not null"`              // 老师：多个老师用逗号分隔
	MaxStudent    int          `json:"max_student" gorm:"not null"`           // 最大选课人数
	WeekHour      int          `json:"week_hour" gorm:"not null"`             // 周学时
	Year          int          `json:"year" gorm:"not null"`                  // 学年
	Semester      int          `json:"semester" gorm:"not null"`              // 学期
	CourseGroupID int          `json:"course_group_id" gorm:"not null;index"` // 课程组类型
	CourseGroup   *CourseGroup `json:"course_group"`
	ReviewCount   int          `json:"review_count" gorm:"not null;default:0"` // 评教数量
	Reviews       ReviewList   `json:"-"`                                      // 所有评教
}

func (c *Course) Create() (err error) {
	err = DB.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Omit(clause.Associations).Create(c).Error
		if err != nil {
			return err
		}

		updateColumns := map[string]any{
			"course_count": gorm.Expr("course_count + 1"),
		}

		// 如果课程组中没有该学分，则添加
		if !slices.Contains(c.CourseGroup.Credits, c.Credit) {
			// 添加学分
			c.CourseGroup.Credits = append(c.CourseGroup.Credits, c.Credit)

			// 手动拼接学分字符串
			var creditsString strings.Builder
			creditsString.WriteByte('[')
			for i, credit := range c.CourseGroup.Credits {
				if i != 0 {
					creditsString.WriteByte(',')
				}
				creditsString.WriteString(strconv.FormatFloat(credit, 'f', -1, 64))
			}
			creditsString.WriteByte(']')
			updateColumns["credits"] = creditsString.String()
		}

		return tx.Model(&CourseGroup{ID: c.CourseGroupID}).
			Updates(updateColumns).Error
	})
	if err != nil {
		return err
	}
	// clear cache
	return common.Cache.Delete(context.Background(), "danke:course_group")
}

type CourseList []*Course

func (l CourseList) LoadReviewList(tx *gorm.DB, options ...FindReviewOption) (err error) {
	var option FindReviewOption
	if len(options) > 0 {
		option = options[0]
	}

	querySet := option.setQuery(tx)

	courseIDs := make([]int, len(l))
	for i, course := range l {
		courseIDs[i] = course.ID
	}

	// 获取课程组的所有课程的所有评论，同时加载评论的历史记录和用户成就
	var reviews ReviewList
	err = querySet.Where("course_id IN ?", courseIDs).Find(&reviews).Error
	if err != nil {
		return err
	}

	// 将评论按照课程分组
	for _, review := range reviews {
		for _, course := range l {
			if course.ID == review.CourseID {
				course.Reviews = append(course.Reviews, review)
			}
		}
	}

	return
}

func (l CourseList) AllReviewList() (reviews ReviewList) {
	for _, course := range l {
		reviews = append(reviews, course.Reviews...)
	}
	return
}
