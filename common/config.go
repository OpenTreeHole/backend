package common

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

// common config
const (
	EnvMode                 = "MODE"
	EnvLogLevel             = "LOG_LEVEL"
	EnvDBType               = "DB_TYPE"
	EnvDBUrl                = "DB_URL"
	EnvCacheType            = "CACHE_TYPE"
	EnvCacheUrl             = "CACHE_URL"
	EnvYiDunBusinessIdText  = "YI_DUN_BUSINESS_ID_TEXT"
	EnvYiDunBusinessIdImage = "YI_DUN_BUSINESS_ID_IMAGE"
	EnvYiDunSecretId        = "YI_DUN_SECRET_ID"
	EnvYiDunSecretKey       = "YI_DUN_SECRET_KEY"
	EnvValidImageUrl        = "VALID_IMAGE_URL"
	EnvUrlHostnameWhitelist = "URL_HOSTNAME_WHITELIST"
	EnvExternalImageHost    = "EXTERNAL_IMAGE_HOST"
	EnvProxyUrl             = "PROXY_URL"
	EnvYiDunAccessKeyId     = "YI_DUN_ACCESS_KEY_ID"
	EnvYiDunAccessKeySecret = "YI_DUN_ACCESS_KEY_SECRET"
)

var defaultConfig = map[string]string{
	EnvMode:                 "dev",
	EnvLogLevel:             "debug",
	EnvDBType:               "sqlite",
	EnvDBUrl:                "file::memory:?cache=shared",
	EnvCacheType:            "memory",
	EnvCacheUrl:             "",
	EnvYiDunBusinessIdText:  "",
	EnvYiDunBusinessIdImage: "",
	EnvYiDunSecretId:        "",
	EnvYiDunSecretKey:       "",
	EnvValidImageUrl:        "",
	EnvUrlHostnameWhitelist: "",
	EnvExternalImageHost:    "",
	EnvProxyUrl:             "",
	EnvYiDunAccessKeyId:     "",
	EnvYiDunAccessKeySecret: "",
}

var GormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // 表名使用单数, `User` -> `user`
	},
	DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束，必须手动创建或者在业务逻辑层维护
	Logger: logger.New(
		mLogger{Logger},
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	),
}

func init() {
	viper.AutomaticEnv()
	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}
