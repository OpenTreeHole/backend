package repository

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
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

type Repository interface {
	Transaction(ctx context.Context, fn func(context.Context) error) error
	GetDB(ctx context.Context) *gorm.DB
	GetCache(ctx context.Context) *cache.Cache
	GetConfig(ctx context.Context) *config.AllConfig
	GetLogger(ctx context.Context) *log.Logger
	// GetFiberCtx 获取 fiber.Ctx
	// 必须在 ctx 中设置 "FiberCtx" 为 fiber.Ctx，否则会 panic
	GetFiberCtx(ctx context.Context) *fiber.Ctx
}

type repository struct {
	db     *gorm.DB
	cache  *cache.Cache
	logger *log.Logger
	conf   *config.AtomicAllConfig
}

func (r *repository) GetLogger(ctx context.Context) *log.Logger {
	//TODO implement me
	return r.logger
}

func NewRepository(db *gorm.DB, cache *cache.Cache, logger *log.Logger, conf *config.AtomicAllConfig) Repository {
	dbMigration(db, conf, logger)
	return &repository{db: db, cache: cache, logger: logger, conf: conf}
}

func NewDB(conf *config.AtomicAllConfig, logger *log.Logger) (db *gorm.DB) {
	var gormConfig = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名使用单数, `User` -> `user`
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束，必须手动创建或者在业务逻辑层维护
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

	var dbConf = conf.Load().DB

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
	dbDsn := conf.Load().DB.DSN
	dbReplicas := conf.Load().DB.Replicas

	if dbDsn == "" {
		logger.Fatal("db url not set")
	}
	switch dbConf.Type {
	case "mysql":
		db = mysqlDB(dbDsn, dbReplicas...)

	case "sqlite":
		db = sqliteDB(dbDsn)

	case "postgres":
		db = postgresDB(dbDsn, dbReplicas...)

	case "memory":
		db = memoryDB()

	default:
		logger.Fatal("db type not support")
	}

	if conf.Load().Mode == "dev" {
		db = db.Debug()
	}
	return
}

func dbMigration(db *gorm.DB, conf *config.AtomicAllConfig, logger *log.Logger) {
	var models = []any{
		model.User{},
		model.UserAchievement{},
		model.Achievement{},
		model.DeleteIdentifier{},
	}
	if conf.Load().Modules.Treehole {
		models = append(models, model.Division{})
	}
	if conf.Load().Modules.CurriculumBoard {
		models = append(
			models,
			model.Review{},
			model.Course{},
			model.CourseGroup{},
			model.ReviewHistory{},
			model.ReviewVote{},
		)
	}

	err := db.AutoMigrate(models...)
	if err != nil {
		logger.Fatal("auto migrate error", zap.Error(err))
	}
}

// Transaction wraps the given function in a transaction.
func (r *repository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return r.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		newCtx := context.WithValue(ctx, "DB", tx)
		return fn(newCtx)
	})
}

func (r *repository) GetDB(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value("DB").(*gorm.DB); ok {
		// check if db is in transaction
		if _, ok := db.Statement.ConnPool.(gorm.TxCommitter); ok {
			return db
		} else {
			return db.WithContext(ctx)
		}
	}

	return r.db.WithContext(ctx)
}

func (r *repository) GetCache(_ context.Context) *cache.Cache {
	return r.cache
}

func (r *repository) GetConfig(_ context.Context) *config.AllConfig {
	return r.conf.Load()
}

func (r *repository) GetFiberCtx(ctx context.Context) *fiber.Ctx {
	return ctx.Value("FiberCtx").(*fiber.Ctx)
}
