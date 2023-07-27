package database

import (
	"context"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/opentreehole/backend/internal/pkg/config"
	"github.com/opentreehole/backend/internal/pkg/database/cache"
)

var Cache cache.Cacher

func initCache() {
	switch config.Config.Cache.Type {
	case "redis":
		Cache = cache.NewRedisCacher(&redis.Options{
			Addr:     config.Config.Cache.Url,
			Username: config.Config.Cache.Username,
			Password: config.Config.Cache.Username,
			DB:       config.Config.Cache.DB,
		})
	case "memory":
		Cache = cache.NewMemoryCacher()
	default:
		log.Fatal().Msgf("unknown cache type: %s", config.Config.Cache.Type)
	}

	log.Debug().Str("type", config.Config.Cache.Type).Msg("init cache success")
}

func GetModelFromCache(key string, model any) bool {
	data, err := Cache.Get(context.Background(), key)
	if err != nil {
		return false
	}
	err = json.Unmarshal(data, model)
	return err == nil
}

func SetModelIntoCache(key string, model any, expiration time.Duration) error {
	data, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return Cache.Set(context.Background(), key, data, expiration)
}

func DeleteModelFromCache(key string) error {
	return Cache.Del(context.Background(), key)
}
