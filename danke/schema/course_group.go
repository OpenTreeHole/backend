package schema

import (
	"github.com/jinzhu/copier"
	"github.com/opentreehole/backend/common"

	"github.com/opentreehole/backend/danke/model"
)

// CourseGroupV1Response 旧版本课程组响应
type CourseGroupV1Response struct {
	ID         int                 `json:"id"`                    // 课程组 ID
	Name       string              `json:"name"`                  // 课程组名称
	Code       string              `json:"code"`                  // 课程组编号
	Department string              `json:"department"`            // 开课学院
	CampusName string              `json:"campus_name"`           // 开课校区
	CourseList []*CourseV1Response `json:"course_list,omitempty"` // 课程组下的课程，slices 必须非空
}

func (r *CourseGroupV1Response) FromModel(
	user *common.User,
	group *model.CourseGroup,
) *CourseGroupV1Response {
	err := copier.Copy(r, group)
	if err != nil {
		panic(err)
	}

	if group.Courses == nil {
		return r
	}
	r.CourseList = make([]*CourseV1Response, 0, len(group.Courses))
	for _, course := range group.Courses {
		r.CourseList = append(r.CourseList, new(CourseV1Response).FromModel(user, course))
	}

	return r
}

type CourseGroupHashV1Response struct {
	Hash string `json:"hash"`
}

func (r *CourseGroupHashV1Response) FromModel(hash string) *CourseGroupHashV1Response {
	r.Hash = hash
	return r
}

/* V3 */

type CourseGroupSearchV3Request struct {
	Query    string `json:"query" form:"query" query:"query" validate:"required" example:"计算机"`
	Page     int    `json:"page" form:"page" query:"page" validate:"min=0" example:"1"`
	PageSize int    `json:"page_size" form:"page_size" query:"page_size" validate:"min=0,max=100" example:"10"`
}

type CourseGroupV3Response struct {
	ID          int                 `json:"id"`                    // 课程组 ID
	Name        string              `json:"name"`                  // 课程组名称
	Code        string              `json:"code"`                  // 课程组编号
	Credits     []float64           `json:"credits"`               // 学分
	Department  string              `json:"department"`            // 开课学院
	CampusName  string              `json:"campus_name"`           // 开课校区
	CourseCount int                 `json:"course_count"`          // 课程数量
	ReviewCount int                 `json:"review_count"`          // 评价数量
	CourseList  []*CourseV1Response `json:"course_list,omitempty"` // 课程组下的课程，slices 必须非空
}

func (r *CourseGroupV3Response) FromModel(
	user *common.User,
	group *model.CourseGroup,
) *CourseGroupV3Response {
	err := copier.Copy(r, group)
	if err != nil {
		panic(err)
	}

	if group.Courses == nil {
		return r
	}
	r.CourseList = make([]*CourseV1Response, 0, len(group.Courses))
	for _, course := range group.Courses {
		r.CourseList = append(r.CourseList, new(CourseV1Response).FromModel(user, course))
	}

	return r
}
