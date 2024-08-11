package main

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	BatchSize = 1000
)

type Course struct {
	ID            int    `json:"id"`
	Name          string `json:"name" gorm:"not null"`    // 课程名称
	Code          string `json:"code" gorm:"not null"`    // 课程编号
	CodeID        string `json:"code_id" gorm:"not null"` // 选课序号。用于区分同一课程编号的不同平行班
	Teachers      string `json:"teachers" gorm:"not null"`
	CourseGroupID int    `json:"course_group_id" gorm:"not null;index"` // 课程组类型
}

type Teacher struct {
	ID   int
	Name string `gorm:"not null"` // 课程组类型
}

type TeacherCourseLink struct {
	TeacherID     int `gorm:"primaryKey;autoIncrement:false"`
	CourseGroupID int `gorm:"primaryKey;autoIncrement:false"` // 课程组类型
}

func GenerateTeacherTabele(DB *gorm.DB) {
	Init()

	// reader := bufio.NewReader(os.Stdin)

	dataMap := map[string][]int{}

	var offset int = 0
	var queryResult []Course
	for {
		query := DB.Table("course").Limit(BatchSize).Offset(offset)
		query.Find(&queryResult)
		offset += len(queryResult)

		if len(queryResult) == 0 {
			// Finished iterating
			break
		}

		for _, course := range queryResult {
			teacherList := strings.Split(course.Teachers, ",")
			for _, name := range teacherList {
				courseList, found := dataMap[name]
				if found {
					dataMap[name] = append(courseList, course.CourseGroupID)
				} else {
					dataMap[name] = []int{course.CourseGroupID}
				}
			}
		}

		// Avoid insertion failure due to duplication

		fmt.Printf("Handled %d records\n", offset)
		// _, _ = reader.ReadString('\n')
	}

	var teachers []*Teacher
	for k := range dataMap {
		teachers = append(teachers, &Teacher{Name: k})
	}

	DB.Clauses(clause.OnConflict{DoNothing: true}).Table("teacher").Create(teachers)

	var links []*TeacherCourseLink
	for index, teacher := range teachers {
		for _, cid := range dataMap[teacher.Name] {
			links = append(links, &TeacherCourseLink{TeacherID: teacher.ID, CourseGroupID: cid})
		}

		// Submit every 100 teachers to avoid SQL being too long
		if index%100 == 0 {
			fmt.Printf("Inserted %d teachers\n", index)

			// Avoid insertion failure due to duplication
			DB.Clauses(clause.OnConflict{DoNothing: true}).Table("teacher_courses").Create(links)
			links = nil
		}
	}
}
