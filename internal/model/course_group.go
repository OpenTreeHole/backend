package model

// CourseGroup 课程组
type CourseGroup struct {
	// 课程组 ID , primary key
	ID int `json:"id"`

	// 课程组名称
	Name string `json:"name" gorm:"not null"`

	// 课程组编号
	Code string `json:"code" gorm:"not null"`

	// 开课学院
	Department string `json:"department" gorm:"not null"`

	// 开课校区
	CampusName string `json:"campus_name" gorm:"not null"`

	// 所有课程
	Courses []*Course `json:"courses"`
}
