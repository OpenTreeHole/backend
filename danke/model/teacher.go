package model

type Teacher struct {
	ID           int            `json:"id"`
	Name         string         `json:"name" gorm:"not null;uniqueIndex:uk_teacher_name,length:191"` // 教师姓名
	CourseGroups []*CourseGroup `gorm:"many2many:teacher_course_groups;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}