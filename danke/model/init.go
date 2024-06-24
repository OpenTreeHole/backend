package model

import (
	"github.com/glebarez/sqlite"
	"github.com/opentreehole/backend/common"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
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

	if viper.GetString(common.EnvLogLevel) == "debug" {
		DB = DB.Debug()
	}

	err = DB.AutoMigrate(
		&CourseGroup{},
		&Course{},
		&Review{},
		&ReviewHistory{},
		&Achievement{},
		&UrlHostnameWhitelist{},
	)
	if err != nil {
		panic(err)
	}
	var hostnames []string
	err = DB.Model(&UrlHostnameWhitelist{}).Pluck("hostname", &hostnames).Error
	if err != nil {
		panic(err)
	}
	viper.Set(common.EnvUrlHostnameWhitelist, hostnames)
}
