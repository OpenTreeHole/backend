package common

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/metrics"
	"github.com/eko/gocache/lib/v4/store"
	bigcache_store "github.com/eko/gocache/store/bigcache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

var Cache struct {
	*marshaler.Marshaler
}

func InitCache() {
	var (
		cacheStore store.StoreInterface
	)

	cacheType := viper.GetString(EnvCacheType)
	cacheUrl := viper.GetString(EnvCacheUrl)
	switch cacheType {
	case "memory":
		bigcacheClient, err := bigcache.New(context.Background(), bigcache.Config{
			Shards:             1024,
			LifeWindow:         10 * time.Minute,
			CleanWindow:        1 * time.Second,
			MaxEntriesInWindow: 1000 * 10 * 60,
			MaxEntrySize:       500,
			StatsEnabled:       false,
			Verbose:            true,
			HardMaxCacheSize:   0,
			Logger:             mLogger{Logger},
		})
		if err != nil {
			panic(err)
		}
		cacheStore = bigcache_store.NewBigcache(bigcacheClient)
	case "redis":
		redisClient := redis.NewClient(&redis.Options{
			Addr: cacheUrl,
		})
		cacheStore = redis_store.NewRedis(redisClient)
	}

	metricsCache := cache.NewMetric[any](
		metrics.NewPrometheus("cache"),
		cache.New[any](cacheStore),
	)

	Cache.Marshaler = marshaler.New(metricsCache)
}
