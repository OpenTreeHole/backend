package model

import (
	. "github.com/opentreehole/backend/common"
	. "github.com/opentreehole/backend/image_hosting/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Hostname string

func Init() {
	var err error
	source := mysql.Open(Config.DbURL)
	DB, err = gorm.Open(source, GormConfig)
	Hostname = "localhost:8000"

	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&ImageTable{})
	if err != nil {
		panic(err)
	}

}
