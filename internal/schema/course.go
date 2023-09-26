package schema

import (
	"github.com/jinzhu/copier"

	"github.com/opentreehole/backend/internal/model"
)

type CourseV1Response struct {
	// 课程 ID
	ID int `json:"id"`

	// 课程名称
	Name string `json:"name"`

	// 课程编号
	Code string `json:"code"`

	// 选课序号。用于区分同一课程编号的不同平行班
	CodeID string `json:"code_id"`

	// 学分
	Credit float64 `json:"credit"`

	// 开课学院
	Department string `json:"department"`

	// 开课校区
	CampusName string `json:"campus_name"`

	// 老师：多个老师用逗号分隔
	Teachers string `json:"teachers"`

	// 最大选课人数
	MaxStudent int `json:"max_student"`

	// 周学时
	WeekHour int `json:"week_hour"`

	// 学年
	Year int `json:"year"`

	// 学期
	Semester int `json:"semester"`

	// 评教列表
	ReviewList []*ReviewV1Response `json:"review_list,omitempty"`
}

func (r *CourseV1Response) FromModel(
	user *model.User,
	course *model.Course,
	votesMap map[int]*model.ReviewVote,
) *CourseV1Response {
	err := copier.Copy(r, course)
	if err != nil {
		panic(err)
	}

	if course.Reviews == nil {
		return r
	}
	r.ReviewList = make([]*ReviewV1Response, 0, len(course.Reviews))
	for _, review := range course.Reviews {
		r.ReviewList = append(r.ReviewList, new(ReviewV1Response).FromModel(user, review, votesMap))
	}

	return r
}

type CreateCourseV1Request struct {
	Name       string  `json:"name" validate:"required,min=1,max=255"`
	Code       string  `json:"code" validate:"required,min=4"`
	CodeID     string  `json:"code_id" validate:"required,min=4"`
	Credit     float64 `json:"credit" validate:"required,min=0.5"`
	Department string  `json:"department" validate:"required,min=1"`
	CampusName string  `json:"campus_name" validate:"required,min=1"`
	Teachers   string  `json:"teachers" validate:"required,min=1"`
	MaxStudent int     `json:"max_student" validate:"required"`
	WeekHour   int     `json:"week_hour" validate:"required"`
	Year       int     `json:"year" validate:"required,min=2000"`
	Semester   int     `json:"semester" validate:"required,min=1"`
}

func (r *CreateCourseV1Request) ToModel(groupID int) *model.Course {
	var course model.Course
	err := copier.Copy(&course, r)
	if err != nil {
		panic(err)
	}
	course.CourseGroupID = groupID
	return &course
}

func (r *CreateCourseV1Request) ToCourseGroupModel() *model.CourseGroup {
	var courseGroup model.CourseGroup
	err := copier.Copy(&courseGroup, r)
	if err != nil {
		panic(err)
	}
	return &courseGroup
}
