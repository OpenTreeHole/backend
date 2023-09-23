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
	Credit int `json:"credit"`

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
	ReviewList []*ReviewV1Response `json:"review_list"`
}

func (r *CourseV1Response) FromModel(
	user *model.User,
	course *model.Course,
) *CourseV1Response {
	err := copier.Copy(r, course)
	if err != nil {
		panic(err)
	}

	for _, review := range course.Reviews {
		r.ReviewList = append(r.ReviewList, new(ReviewV1Response).FromModel(user, review))
	}

	return r
}
