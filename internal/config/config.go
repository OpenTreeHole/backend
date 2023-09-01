package config

import (
	"github.com/creasty/defaults"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var Config *Type

type Type struct {
	// app mode: dev, production, test, bench, default is dev
	Mode string `yaml:"mode" default:"dev"`

	// if true, log above debug level, else log above info level
	Debug bool `yaml:"debug" default:"false"`

	// set port, default 8000
	Port int `yaml:"port"`

	// relational database settings
	DB struct {
		// mysql or sqlite or postgres or memory, case-insensitive, default is sqlite for dev
		Type string `yaml:"type" default:"sqlite"`

		// DSN is the data source name
		//
		// mysql example: user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=true&loc=Asia%2fShanghai
		// set time_zone in url, otherwise UTC
		// for more detail, see https://github.com/go-sql-driver/mysql#dsn-data-source-name
		//
		// memory example: file::memory:?cache=shared
		// for more detail, see https://gorm.io/docs/connecting_to_the_database.html#SQLite
		//
		// sqlite example: data/sqlite.db
		// for more detail, see https://gorm.io/docs/connecting_to_the_database.html#SQLite
		//
		// postgres example: host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
		// for more detail, see https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
		DSN string `yaml:"dsn" default:"data/sqlite.db"`

		// when type is mysql or postgres, set this to enable read-write separation
		Replicas []string `yaml:"replicas"`
	} `yaml:"db,omitempty"`

	// database cache settings
	Cache struct {
		// redis or memory, case-insensitive, default is memory for dev
		Type string `yaml:"type" default:"memory"`

		// redis url
		Url string `yaml:"url"`

		// redis username
		Username string `yaml:"username"`

		// redis password
		Password string `yaml:"password"`

		// redis db number, default is 0
		DB int `yaml:"db"`
	} `yaml:"cache,omitempty"`

	// search engine settings
	SearchEngine struct {
		// Elasticsearch or Meilisearch, case-insensitive, default is Elasticsearch
		Type string `yaml:"type" default:"elasticsearch"`

		// search engine url
		Url string `yaml:"url"`
	} `yaml:"searchEngine,omitempty"`

	Gateway struct {
		// api gateway type
		// support: kong, apisix
		Type string `yaml:"type" default:"kong"`

		// api gateway url
		Url string `yaml:"url" default:"kong:8001"`

		// api gateway token
		Token string `yaml:"token"`
	} `yaml:"gateway,omitempty"`

	// smtp server settings, used for email verification and notification
	Email struct {
		// smtp server address
		Host string `yaml:"host"`

		// smtp server port
		Port int `yaml:"port" default:"465"`

		// smtp server username
		Username string `yaml:"username"`

		// smtp server password
		Password string `yaml:"password"`

		// smtp server from address
		From string `yaml:"from"`

		// smtp server from name
		FromName string `yaml:"fromName"`

		// smtp server tls
		TLS bool `yaml:"tls" default:"true"`

		// email white list
		WhiteList []string `yaml:"whiteList"`

		// site name, using when send email
		SiteName string `yaml:"siteName" default:"Open Tree Hole"`

		// dev email, using when debug
		DevEmail string `yaml:"devEmail"`
	} `yaml:"email,omitempty"`

	// feature settings, including standalone, shamir, emailVerification, emailNotification, registrationTest
	Feature struct {
		// enable standalone mode, means jwt-auth without api gateway
		Standalone bool `yaml:"standalone"`

		// enable shamir secret sharing encryption for email
		Shamir bool `yaml:"shamir" default:"true"`

		// enable email verification, default is false
		EmailVerification bool `yaml:"emailVerification"`

		VerificationCodeExpires int `yaml:"verificationCodeExpires" default:"600"` // seconds

		// enable email notification, default is false
		EmailNotification bool `yaml:"emailNotification"`

		// enable registration test, default is false
		RegistrationTest bool `yaml:"registrationTest"`
	} `yaml:"feature,omitempty"`

	// notification settings, used for notification server
	Notification struct {
		// notification certificates and package name
		MipushKeyPath      string `yaml:"mipushKeyPath" default:"data/mipush.pem"`
		APNSKeyPath        string `yaml:"apnsKeyPath" default:"data/apns.pem"`
		IOSPackageName     string `yaml:"iosPackageName" default:"io.github.danxi-dev.dan-xi"`
		AndroidPackageName string `yaml:"androidPackageName" default:"io.github.danxi_dev.dan_xi"`

		// mipush notification callback url, used for notification server
		MipushCallbackUrl string `yaml:"mipushCallbackUrl" default:"http://notification.fduhole.com/api/callback/mipush"`
	}
}

func InitConfig() (conf *viper.Viper) {
	var configFile = "config/local.yml"
	pflag.StringP("config", "c", configFile, "config file path")

	conf = viper.New()
	err := viper.BindPFlag("configPath", pflag.Lookup("config"))
	if err != nil {
		panic(err)
		//zap.L().Fatal("failed to bind config flag", zap.Error(err))
	}
	conf.SetConfigFile(configFile)
	err = conf.ReadInConfig()
	if err != nil {
		panic(err)
		//zap.L().Fatal("failed to load config file", zap.Error(err))
		return nil
	}

	// parse config from file
	err = conf.Unmarshal(&Config)
	if err != nil {
		panic(err)
		//zap.L().Fatal("failed to parse config file", zap.Error(err))
		return nil
	}

	// set default value
	_ = defaults.Set(&Config)

	return
}
