package config

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"sync/atomic"

	"github.com/creasty/defaults"
	"github.com/spf13/pflag"

	"github.com/opentreehole/backend/pkg/utils"
)

type Config struct {
	// app mode: dev, production, test, bench, default is dev
	Mode string `yaml:"mode" default:"dev" json:"mode"`

	// if true, log above debug level, else log above info level
	Debug bool `yaml:"debug" default:"false" json:"debug"`

	// set port, default 8000
	Port int `yaml:"port" json:"port"`

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
		Url string `yaml:"url" json:"url"`

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
		Url string `yaml:"url" json:"url"`
	} `yaml:"searchEngine" json:"search_engine"`

	// api gateway settings
	// if standalone is true, gateway settings will be ignored
	Gateway struct {
		// api gateway type
		// support: kong, apisix
		Type string `yaml:"type" default:"kong" json:"type"`

		// api gateway url
		Url string `yaml:"url" default:"kong:8001" json:"url"`

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
		FromName string `yaml:"fromName" json:"from_name"`

		// smtp server tls
		TLS bool `yaml:"tls" default:"true" json:"tls"`

		// email white list
		WhiteList []string `yaml:"whiteList" json:"white_list"`

		// site name, using when send email
		SiteName string `yaml:"siteName" default:"Open Tree Hole" json:"site_name"`

		// dev email, using when debug
		DevEmail string `yaml:"devEmail" json:"dev_email"`
	} `yaml:"email" json:"email"`

	// feature settings, including standalone, shamir, emailVerification, emailNotification, registrationTest
	Feature struct {
		// enable standalone mode, means jwt-auth without api gateway, default is true
		Standalone bool `yaml:"standalone" default:"true" json:"standalone"`

		// enable shamir secret sharing encryption for email, default is false
		Shamir bool `yaml:"shamir" default:"false" json:"shamir"`

		// enable email verification, default is false
		EmailVerification bool `yaml:"emailVerification" default:"false" json:"email_verification"`

		// email verification code expires, default is 600 seconds
		VerificationCodeExpires int `yaml:"verificationCodeExpires" default:"600" json:"verification_code_expires"`

		// enable email notification, default is false
		EmailNotification bool `yaml:"emailNotification" default:"false" json:"email_notification"`

		// enable registration test, default is false
		RegistrationTest bool `yaml:"registrationTest" default:"false" json:"registration_test"`
	} `yaml:"feature" json:"feature"`

	// notification settings, used for notification server
	Notification struct {
		// notification certificates and package name
		MipushKeyPath      string `yaml:"mipushKeyPath" default:"data/mipush.pem" json:"mipush_key_path"`
		APNSKeyPath        string `yaml:"apnsKeyPath" default:"data/apns.pem" json:"apns_key_path"`
		IOSPackageName     string `yaml:"iosPackageName" default:"io.github.danxi-dev.dan-xi" json:"ios_package_name"`
		AndroidPackageName string `yaml:"androidPackageName" default:"io.github.danxi_dev.dan_xi" json:"android_package_name"`

		// mipush notification callback url, used for notification server
		MipushCallbackUrl string `yaml:"mipushCallbackUrl" default:"http://notification.fduhole.com/api/callback/mipush" json:"mipush_callback_url"`
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

	// read config from file
	err = config.ReadFromFile(configFilename)
	if err != nil {
		panic(err)
	}

	// save config
	err = config.WriteIntoFile(configFilename)
	if err != nil {
		panic(err)
	}

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

	allConfig.Store(&AllConfig{
		Config:     &config,
		FileConfig: &fileConfig,
	})
	return &allConfig
}

// ReadFromFile read config from file
// if file not exist, create it with default value; else read it
func (config *Config) ReadFromFile(name string) (err error) {
	var file *os.File

	// set default value finally
	defer defaults.MustSet(config)

	file, err = os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return
	}

	return
}

// WriteIntoFile write config into file
// if file not exist, create it; else truncate it
func (config *Config) WriteIntoFile(name string) (err error) {
	var file *os.File

	file, err = os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return
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
		return
	}

	return
}
