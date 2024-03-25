package model

import (
	"github.com/opentreehole/backend/common"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error

	dbType := viper.GetString(common.EnvDBType)
	dbUrl := viper.GetString(common.EnvDBUrl)
	switch dbType {
	case "sqlite":
		if dbUrl == "" {
			dbUrl = "sqlite.db"
		}
		DB, err = gorm.Open(sqlite.Open(dbUrl), common.GormConfig)
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dbUrl), common.GormConfig)
	case "postgres":
		DB, err = gorm.Open(postgres.Open(dbUrl), common.GormConfig)
	default:
		panic("db type not support")
	}
	if err != nil {
		panic(err)
	}

	if viper.GetString(common.EnvMode) == "dev" {
		DB = DB.Debug()
	}

	err = DB.AutoMigrate(
		&CourseGroup{},
		&Course{},
		&Review{},
		&Achievement{},
	)
	if err != nil {
		panic(err)
	}
}
