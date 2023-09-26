package config

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"sync/atomic"

	"github.com/caarlos0/env/v9"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/pflag"

	"github.com/opentreehole/backend/pkg/utils"
)

type EnvConfig struct {
	Mode string `env:"MODE" default:"dev" validate:"oneof=dev production test bench"`

	LogLevel string `env:"LOG_LEVEL" default:"debug" validate:"oneof=debug info warn error dpanic panic fatal"`

	Port int `env:"PORT" default:"8000"`

	DBType string `env:"DB_TYPE" default:"sqlite" validate:"oneof=mysql sqlite postgres memory"`

	DBDSN string `env:"DB_DSN" default:"data/sqlite.db"`

	CacheType string `env:"CACHE_TYPE" default:"memory" validate:"oneof=redis memory"`

	CacheUrl string `env:"CACHE_URL" default:"redis:6379"`

	SearchEngineType string `env:"SEARCH_ENGINE_TYPE" default:"elasticsearch" validate:"oneof=elasticsearch meilisearch"`

	SearchEngineUrl string `env:"SEARCH_ENGINE_URL" default:"http://elasticsearch:9200"`

	ModulesAuth bool `env:"MODULES_AUTH" default:"false"`

	ModulesNotification bool `env:"MODULES_NOTIFICATION" default:"false"`

	ModulesTreehole bool `env:"MODULES_TREEHOLE" default:"false"`

	ModulesCurriculumBoard bool `env:"MODULES_CURRICULUM_BOARD" default:"false"`
}

type Config struct {
	// app mode: dev, production, test, bench, default is dev
	Mode string `yaml:"mode" default:"dev" json:"mode" validate:"oneof=dev production test bench"`

	// LogLevel is the log level, default is debug
	LogLevel string `yaml:"log_level" default:"debug" json:"log_level" validate:"oneof=debug info warn error dpanic panic fatal"`

	// set port, default 8000
	Port int `yaml:"port" json:"port" default:"8000"`

	// relational database settings
	DB struct {
		// mysql or sqlite or postgres or memory, case-insensitive, default is sqlite for dev
		Type string `yaml:"type" default:"sqlite" json:"type"`

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
		DSN string `yaml:"dsn" default:"data/sqlite.db" json:"dsn"`

		// when type is mysql or postgres, set this to enable read-write separation
		Replicas []string `yaml:"replicas" json:"replicas"`
	} `yaml:"db" json:"db"`

	// database cache settings
	Cache struct {
		// redis or memory, case-insensitive, default is memory for dev
		Type string `yaml:"type" default:"memory" json:"type"`

		// redis url
		Url string `yaml:"url" default:"redis:6379" json:"url"`

		// redis username
		Username string `yaml:"username" json:"username"`

		// redis password
		Password string `yaml:"password" json:"password"`

		// redis db number, default is 0
		DB int `yaml:"db" json:"db"`
	} `yaml:"cache" json:"cache"`

	// search engine settings
	SearchEngine struct {
		// Elasticsearch or Meilisearch, case-insensitive, default is Elasticsearch
		Type string `yaml:"type" default:"elasticsearch" json:"type"`

		// search engine url
		Url string `yaml:"url" default:"http://elasticsearch:9200" json:"url"`
	} `yaml:"searchEngine" json:"search_engine"`

	// api gateway settings
	// if standalone is true, gateway settings will be ignored
	Gateway struct {
		// api gateway type
		// support: kong, apisix
		Type string `yaml:"type" default:"kong" json:"type"`

		// api gateway url
		Url string `yaml:"url" default:"http://kong:8001" json:"url"`

		// api gateway token
		Token string `yaml:"token" json:"token"`
	} `yaml:"gateway" json:"gateway"`

	// smtp server settings, used for email verification and notification
	Email struct {
		// smtp server address
		Host string `yaml:"host" json:"host"`

		// smtp server port
		Port int `yaml:"port" default:"465" json:"port"`

		// smtp server username
		Username string `yaml:"username" json:"username"`

		// smtp server password
		Password string `yaml:"password" json:"password"`

		// smtp server from address
		From string `yaml:"from" json:"from"`

		// smtp server from name
		FromName string `yaml:"from_name" json:"from_name"`

		// smtp server tls
		TLS bool `yaml:"tls" default:"true" json:"tls"`

		// email white list
		WhiteList []string `yaml:"white_list" json:"white_list"`

		// site name, using when send email
		SiteName string `yaml:"site_name" default:"Open Tree Hole" json:"site_name"`

		// dev email, using when debug
		DevEmail string `yaml:"dev_email" json:"dev_email"`
	} `yaml:"email" json:"email"`

	Modules struct {
		// enable auth module
		Auth bool `yaml:"auth" default:"false" json:"auth"`

		// enable notification module
		Notification bool `yaml:"notification" default:"false" json:"notification"`

		// enable treehole module
		Treehole bool `yaml:"treehole" default:"false" json:"treehole"`

		// enable curriculum_board module
		CurriculumBoard bool `yaml:"curriculum_board" default:"false" json:"curriculum_board"`
	} `yaml:"modules" json:"modules"`

	// feature settings
	Features struct {
		// enable external gateway mode, means jwt-auth with api gateway
		ExternalGateway bool `yaml:"external_gateway" default:"false" json:"external_gateway"`

		// enable shamir secret sharing encryption for email
		Shamir bool `yaml:"shamir" default:"false" json:"shamir"`

		// enable email verification
		EmailVerification bool `yaml:"email_verification" default:"false" json:"email_verification"`

		// email verification code expires, default is 600 seconds
		VerificationCodeExpires int `yaml:"verification_code_expires" default:"600" json:"verification_code_expires"`

		// enable email notification
		EmailNotification bool `yaml:"email_notification" default:"false" json:"email_notification"`

		// enable registration test
		RegistrationTest bool `yaml:"registration_test" default:"false" json:"registration_test"`
	} `yaml:"features" json:"features"`

	// notification settings, used for notification server
	Notification struct {
		// notification certificates and package name
		MipushKeyPath      string `yaml:"mipush_key_path" default:"data/mipush.pem" json:"mipush_key_path"`
		APNSKeyPath        string `yaml:"apns_key_path" default:"data/apns.pem" json:"apns_key_path"`
		IOSPackageName     string `yaml:"ios_package_name" default:"io.github.danxi-dev.dan-xi" json:"ios_package_name"`
		AndroidPackageName string `yaml:"android_package_name" default:"io.github.danxi_dev.dan_xi" json:"android_package_name"`

		// mipush notification callback url, used for notification server
		MipushCallbackUrl string `yaml:"mipush_callback_url" default:"http://notification.fduhole.com/api/callback/mipush" json:"mipush_callback_url"`
	} `json:"notification"`
}

type FileConfig struct {
	// DecryptedIdentifierSalt is the decrypted identifier salt
	DecryptedIdentifierSalt []byte
}

type AllConfig struct {
	*Config
	*FileConfig
}

type AtomicAllConfig = atomic.Pointer[AllConfig]

var defaultConfig Config

func init() {
	defaults.MustSet(&defaultConfig)
}

func NewConfig() *AtomicAllConfig {
	var (
		err                    error
		configFilename         string
		identifierSaltFilename string
		config                 Config
		fileConfig             FileConfig
		allConfig              AtomicAllConfig
	)

	const (
		defaultConfigFile         = "config/config.json"
		defaultIdentifierSaltFile = "data/identifier_salt"
	)

	pflag.StringVarP(
		&configFilename,
		"config",
		"c",
		defaultConfigFile,
		"config file path",
	)
	pflag.StringVarP(
		&identifierSaltFilename,
		"identifierSalt",
		"s",
		defaultIdentifierSaltFile,
		"identifier salt file path",
	)
	pflag.Parse()

	// get env config
	envConfig := GetEnvConfig()

	// read config from file
	config.ReadFromFile(configFilename)

	// copy env config to config
	CopyEnvConfigToConfig(envConfig, &config)

	err = validator.New().Struct(&config)
	if err != nil {
		panic(err)
	}

	// save config
	config.WriteIntoFile(configFilename)

	if config.Modules.Auth {
		// parse identifier salt from file
		identifierSaltBytes, err := os.ReadFile(identifierSaltFilename)
		if err != nil {
			if os.IsNotExist(err) {
				if config.Mode == "production" {
					panic("identifier salt file not found")
				} else {
					fileConfig.DecryptedIdentifierSalt = []byte("123456")
				}
			}
		} else {
			fileConfig.DecryptedIdentifierSalt, err = base64.StdEncoding.DecodeString(utils.BytesToString(identifierSaltBytes))
			if err != nil {
				panic("decode identifier salt error")
			}
		}
	}

	allConfig.Store(&AllConfig{
		Config:     &config,
		FileConfig: &fileConfig,
	})
	return &allConfig
}

func GetEnvConfig() *EnvConfig {
	var envConfig EnvConfig
	err := env.Parse(&envConfig)
	if err != nil {
		panic(err)
	}

	defer defaults.MustSet(&envConfig)

	return &envConfig
}

// ReadFromFile read config from file
// if file not exist, create it with default value; else read it
func (config *Config) ReadFromFile(name string) {
	var file *os.File

	// set default value finally
	defer defaults.MustSet(config)

	file, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		panic(err)
	}
}

// WriteIntoFile write config into file
// if file not exist, create it; else truncate it
func (config *Config) WriteIntoFile(name string) {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(config)
	if err != nil {
		panic(err)
	}

}

func CopyEnvConfigToConfig(envConfig *EnvConfig, config *Config) {
	config.Mode = envConfig.Mode
	config.LogLevel = envConfig.LogLevel
	config.Port = envConfig.Port
	config.DB.Type = envConfig.DBType
	config.DB.DSN = envConfig.DBDSN
	config.Cache.Type = envConfig.CacheType
	config.Cache.Url = envConfig.CacheUrl
	config.SearchEngine.Type = envConfig.SearchEngineType
	config.SearchEngine.Url = envConfig.SearchEngineUrl
	config.Modules.Auth = envConfig.ModulesAuth
	config.Modules.Notification = envConfig.ModulesNotification
	config.Modules.Treehole = envConfig.ModulesTreehole
	config.Modules.CurriculumBoard = envConfig.ModulesCurriculumBoard
}
