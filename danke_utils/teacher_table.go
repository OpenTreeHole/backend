package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/opentreehole/backend/common"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	BatchSize = 1000
)

var GormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // 表名使用单数, `User` -> `user`
	},
	DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束，必须手动创建或者在业务逻辑层维护
	Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	),
}

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

var DB *gorm.DB

func Init() {
	viper.AutomaticEnv()
	dbType := viper.GetString(common.EnvDBType)
	dbUrl := viper.GetString(common.EnvDBUrl)

	var err error

	switch dbType {
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dbUrl), GormConfig)
	case "postgres":
		DB, err = gorm.Open(postgres.Open(dbUrl), GormConfig)
	default:
		panic("db type not supported")
	}

	if err != nil {
		panic(err)
	}
}

func main() {
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
