package database

import (
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	"github.com/opentreehole/backend/internal/pkg/config"
)

var DB *gorm.DB

var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
	},
	Logger: logger.New(
		&log.Logger,
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	),
}

// Read/Write Splitting
func mysqlDB(dsn string, replicasDsn ...string) *gorm.DB {
	// set source databases
	source := mysql.Open(dsn)
	db, err := gorm.Open(source, gormConfig)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// set replica databases
	if len(replicasDsn) == 0 {
		return db
	}
	var replicas []gorm.Dialector
	for _, url := range config.Config.DB.Replicas {
		replicas = append(replicas, mysql.Open(url))
	}
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{source},
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return db
}

func postgresDB(dsn string, replicasDsn ...string) *gorm.DB {
	// set source databases
	source := postgres.Open(dsn)
	db, err := gorm.Open(source, gormConfig)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	// set replica databases
	if len(replicasDsn) == 0 {
		return db
	}
	var replicas []gorm.Dialector
	for _, url := range config.Config.DB.Replicas {
		replicas = append(replicas, postgres.Open(url))
	}
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{source},
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return db
}

func sqliteDB(filePath string) *gorm.DB {
	// create file if not exist
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		_, err = os.Create(filePath)
		if err != nil {
			log.Trace().Err(err).Send()
		}
	}

	db, err := gorm.Open(sqlite.Open(filePath), gormConfig)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	// https://github.com/go-gorm/gorm/issues/3709
	phyDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	phyDB.SetMaxOpenConns(1)
	return db
}

func memoryDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), gormConfig)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	// https://github.com/go-gorm/gorm/issues/3709
	phyDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	phyDB.SetMaxOpenConns(1)
	return db
}

// initGorm init gorm
func initGorm(models ...any) {

	// init db
	if config.Config.DB.Type == "" {
		config.Config.DB.Type = "sqlite"
	}

	switch config.Config.DB.Type {
	case "mysql":
		if config.Config.DB.DSN == "" {
			log.Fatal().Msg("mysql url not set")
		} else {
			DB = mysqlDB(config.Config.DB.DSN, config.Config.DB.Replicas...)
		}
	case "sqlite":
		if config.Config.DB.DSN == "" {
			config.Config.DB.DSN = "data/sqlite.db"
		}
		DB = sqliteDB(config.Config.DB.DSN)
	case "postgres":
		if config.Config.DB.DSN == "" {
			log.Fatal().Msg("postgres url not set")
		}
		DB = postgresDB(config.Config.DB.DSN, config.Config.DB.Replicas...)
	case "memory":
		DB = memoryDB()
	default:
		log.Fatal().Msg("unsupported db type")
	}

	if config.Config.Debug {
		DB = DB.Debug()
	}

	// migrate models
	if len(models) != 0 {
		err := DB.AutoMigrate(models...)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}
}
