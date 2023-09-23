package model

// Course 课程
type Course struct {
	// 课程 ID , primary key
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

	// 课程组类型
	CourseGroupID int `json:"course_group_id"`
}
