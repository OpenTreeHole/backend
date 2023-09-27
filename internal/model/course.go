package model

import (
	"time"
)

// Course 课程
type Course struct {
	ID            int       `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name" gorm:"not null"`                   // 课程名称
	Code          string    `json:"code" gorm:"not null"`                   // 课程编号
	CodeID        string    `json:"code_id" gorm:"not null"`                // 选课序号。用于区分同一课程编号的不同平行班
	Credit        float64   `json:"credit" gorm:"not null"`                 // 学分
	Department    string    `json:"department" gorm:"not null"`             // 开课学院
	CampusName    string    `json:"campus_name" gorm:"not null"`            // 开课校区
	Teachers      string    `json:"teachers" gorm:"not null"`               // 老师：多个老师用逗号分隔
	MaxStudent    int       `json:"max_student" gorm:"not null"`            // 最大选课人数
	WeekHour      int       `json:"week_hour" gorm:"not null"`              // 周学时
	Year          int       `json:"year" gorm:"not null"`                   // 学年
	Semester      int       `json:"semester" gorm:"not null"`               // 学期
	CourseGroupID int       `json:"course_group_id" gorm:"not null;index"`  // 课程组类型
	ReviewCount   int       `json:"review_count" gorm:"not null;default:0"` // 评教数量
	Reviews       []*Review `json:"reviews"`                                // 所有评教
}
