package model

type Teacher struct {
	ID           int            `json:"id"`
	Name         string         `json:"name" gorm:"not null"` // 教师姓名
	CourseGroups []*CourseGroup `gorm:"many2many:teacher_course_link;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}