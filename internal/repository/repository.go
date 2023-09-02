package repository

import (
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/internal/model"
	"github.com/opentreehole/backend/internal/pkg/cache"
	"github.com/opentreehole/backend/pkg/log"
)

type Repository struct {
	db     *gorm.DB
	cache  *cache.Cache
	logger *log.Logger
	conf   *config.AtomicAllConfig
}

func NewRepository(db *gorm.DB, cache *cache.Cache, logger *log.Logger, conf *config.AtomicAllConfig) *Repository {
	return &Repository{db: db, cache: cache, logger: logger, conf: conf}
}

func NewDB(conf *config.AtomicAllConfig, logger *log.Logger) (db *gorm.DB) {
	var gormConfig = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: gormlogger.New(
			zap.NewStdLog(logger.Logger),
			gormlogger.Config{
				SlowThreshold:             time.Second,      // 慢 SQL 阈值
				LogLevel:                  gormlogger.Error, // 日志级别
				IgnoreRecordNotFoundError: true,             // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,            // 禁用彩色打印
			},
		),
	}

	// Read/Write Splitting
	mysqlDB := func(dsn string, replicasDsn ...string) *gorm.DB {
		// set source databases
		source := mysql.Open(dsn)
		db, err := gorm.Open(source, gormConfig)
		if err != nil {
			logger.Fatal("mysql open error", zap.Error(err))
		}

		// set replica databases
		if len(replicasDsn) == 0 {
			return db
		}
		var replicas []gorm.Dialector
		for _, url := range replicasDsn {
			replicas = append(replicas, mysql.Open(url))
		}
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{source},
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
		if err != nil {
			logger.Fatal("mysql replica error", zap.Error(err))
		}
		return db
	}

	postgresDB := func(dsn string, replicasDsn ...string) *gorm.DB {
		// set source databases
		source := postgres.Open(dsn)
		db, err := gorm.Open(source, gormConfig)
		if err != nil {
			logger.Fatal("postgres open error", zap.Error(err))
		}

		// set replica databases
		if len(replicasDsn) == 0 {
			return db
		}
		var replicas []gorm.Dialector
		for _, url := range replicasDsn {
			replicas = append(replicas, postgres.Open(url))
		}
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{source},
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
		if err != nil {
			logger.Fatal("postgres replica error", zap.Error(err))
		}
		return db
	}

	sqliteDB := func(filePath string) *gorm.DB {
		// create file if not exist
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			_, err = os.Create(filePath)
			if err != nil {
				logger.Fatal("sqlite create error", zap.Error(err))
			}
		}

		db, err := gorm.Open(sqlite.Open(filePath), gormConfig)
		if err != nil {
			logger.Fatal("sqlite open error", zap.Error(err))
		}
		// https://github.com/go-gorm/gorm/issues/3709
		phyDB, err := db.DB()
		if err != nil {
			logger.Fatal("sqlite db error", zap.Error(err))
		}
		phyDB.SetMaxOpenConns(1)
		return db
	}

	memoryDB := func() *gorm.DB {
		db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), gormConfig)
		if err != nil {
			logger.Fatal("memory db open error", zap.Error(err))
		}
		// https://github.com/go-gorm/gorm/issues/3709
		phyDB, err := db.DB()
		if err != nil {
			logger.Fatal("memory db error", zap.Error(err))
		}
		phyDB.SetMaxOpenConns(1)
		return db
	}

	// init db
	dbType := conf.Load().DB.Type
	dbDsn := conf.Load().DB.DSN
	dbReplicas := conf.Load().DB.Replicas

	switch dbType {
	case "mysql":
		if dbDsn == "" {
			logger.Fatal("mysql url not set")
		}
		db = mysqlDB(dbDsn, dbReplicas...)

	case "sqlite":
		if dbDsn == "" {
			logger.Fatal("sqlite url not set")
		}
		db = sqliteDB(dbDsn)

	case "postgres":
		if dbDsn == "" {
			logger.Fatal("postgres url not set")
		}
		db = postgresDB(dbDsn, dbReplicas...)

	case "memory":
		db = memoryDB()

	default:
		logger.Fatal("db type not support")
	}

	if conf.Load().Debug {
		db = db.Debug()
	}

	err := db.AutoMigrate(
		model.User{},
	)
	if err != nil {
		logger.Fatal("auto migrate error", zap.Error(err))
	}

	return
}
