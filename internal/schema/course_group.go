package schema

import (
	"github.com/jinzhu/copier"

	"github.com/opentreehole/backend/internal/model"
)

// CourseGroupV1Response 旧版本课程组响应
type CourseGroupV1Response struct {
	// 课程组 ID
	ID int `json:"id"`

	// 课程组名称
	Name string `json:"name"`

	// 课程组编号
	Code string `json:"code"`

	// 开课学院
	Department string `json:"department"`

	// 开课校区
	CampusName string `json:"campus_name"`

	// 课程组下的课程，slices 必须非空
	CourseList []*CourseV1Response `json:"course_list"`
}

func (r *CourseGroupV1Response) FromModel(
	user *model.User,
	group *model.CourseGroup,
) *CourseGroupV1Response {
	err := copier.Copy(r, group)
	if err != nil {
		panic(err)
	}

	for _, course := range group.Courses {
		r.CourseList = append(r.CourseList, new(CourseV1Response).FromModel(user, course))
	}

	return r
}
