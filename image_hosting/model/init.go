package model

import (
	. "github.com/opentreehole/backend/common"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	dbUrl := viper.GetString(EnvDBUrl)
	source := mysql.Open(dbUrl)
	DB, err = gorm.Open(source, GormConfig)
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&ImageTable{})
	if err != nil {
		panic(err)
	}

}
