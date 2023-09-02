package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/metrics"
	"github.com/eko/gocache/lib/v4/store"
	bigcache_store "github.com/eko/gocache/store/bigcache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/pkg/log"
)

type Cache struct {
	*marshaler.Marshaler
}

func NewCache(conf *config.AtomicAllConfig, logger *log.Logger) *Cache {
	var (
		cacheConf  = conf.Load().Cache
		cacheType  = cacheConf.Type
		cacheStore store.StoreInterface
	)
	if cacheType == "memory" {
		bigcacheClient, err := bigcache.New(context.Background(), bigcache.Config{
			Shards:             1024,
			LifeWindow:         10 * time.Minute,
			CleanWindow:        1 * time.Second,
			MaxEntriesInWindow: 1000 * 10 * 60,
			MaxEntrySize:       500,
			StatsEnabled:       false,
			Verbose:            true,
			HardMaxCacheSize:   0,
			Logger:             zap.NewStdLog(logger.Logger),
		})
		if err != nil {
			panic(err)
		}
		cacheStore = bigcache_store.NewBigcache(bigcacheClient)
	} else if cacheType == "redis" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     cacheConf.Url,
			Username: cacheConf.Username,
			Password: cacheConf.Password,
			DB:       cacheConf.DB,
		})
		cacheStore = redis_store.NewRedis(redisClient)
	} else {
		panic("unsupported cache")
	}

	metricsCache := cache.NewMetric[any](
		metrics.NewPrometheus("cache"),
		cache.New[any](cacheStore),
	)

	marshal := marshaler.New(metricsCache)
	return &Cache{Marshaler: marshal}
}
