package main

import (
	"fmt"
	"slices"
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
	CourseGroupID int    `json:"course_group_id" gorm:"not null;index"` // 课程组编号
}

type Teacher struct {
	ID   int
	Name string `gorm:"not null"` // 课程组 ID
}

type TeacherCourseLink struct {
	TeacherID     int `gorm:"primaryKey;autoIncrement:false"`
	CourseGroupID int `gorm:"primaryKey;autoIncrement:false"` // 课程组编号
}

func AppendUnique[T comparable](slice []T, elems ...T) []T {
	for _, elem := range elems {
		if !slices.Contains(slice, elem) {
			slice = append(slice, elem)
		}
	}

	return slice
}

func GenerateTeacherTable(DB *gorm.DB) {
	Init()

	// reader := bufio.NewReader(os.Stdin)

	dataMap := map[string][]int{}

	var queryResult []Course
	query := DB.Table("course")
	query.FindInBatches(&queryResult, BatchSize, func(tx *gorm.DB, batch int) error {
		for _, course := range queryResult {
			teacherList := strings.Split(course.Teachers, ",")
			for _, name := range teacherList {
				courseList, found := dataMap[name]
				if found {
					dataMap[name] = AppendUnique(courseList, course.CourseGroupID)
				} else {
					dataMap[name] = []int{course.CourseGroupID}
				}
			}
		}

		fmt.Printf("Handled batchg %d\n", batch)
		return nil
		})

	var teachers []*Teacher
	for k := range dataMap {
		teachers = append(teachers, &Teacher{Name: k})
	}

	// Avoid insertion failure due to duplication
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
