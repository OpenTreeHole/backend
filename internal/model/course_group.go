package model

// CourseGroup 课程组
type CourseGroup struct {
	// 课程组 ID , primary key
	ID int `json:"id"`

	// 课程组名称
	Name string `json:"name"`

	// 课程组编号
	Code string `json:"code"`

	// 开课学院
	Department string `json:"department"`

	// 开课校区
	CampusName string `json:"campus_name"`
}
