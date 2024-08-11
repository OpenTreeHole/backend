package main

import (
	"log"
	"os"
	"time"

	"github.com/opentreehole/backend/common"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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
	// Call any script as needed
	// GenerateTeacherTabele(DB)
}