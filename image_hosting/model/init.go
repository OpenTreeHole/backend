package model

import (
	"github.com/opentreehole/backend/common"
	"github.com/opentreehole/backend/image_hosting/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Init() {
	source := mysql.Open(config.Config.DbURL)
	DB, err := gorm.Open(source, common.GormConfig)

	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(
		&ImageTable{})

	if err != nil {
		panic(err)
	}
}
