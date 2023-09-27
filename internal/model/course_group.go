package model

// CourseGroup 课程组
type CourseGroup struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" gorm:"not null"`                   // 课程组名称
	Code        string    `json:"code" gorm:"not null"`                   // 课程组编号
	Department  string    `json:"department" gorm:"not null"`             // 开课学院
	CampusName  string    `json:"campus_name" gorm:"not null"`            // 开课校区
	CourseCount int       `json:"course_count" gorm:"not null;default:0"` // 课程数量
	Courses     []*Course `json:"courses"`
}
