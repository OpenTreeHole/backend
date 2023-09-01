package database

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/internal/pkg/database/cache"
)

func initCache() (Cacher cache.Cacher) {
	switch config.Config.Cache.Type {
	case "redis":
		Cacher = cache.NewRedisCacher(&redis.Options{
			Addr:     config.Config.Cache.Url,
			Username: config.Config.Cache.Username,
			Password: config.Config.Cache.Username,
			DB:       config.Config.Cache.DB,
		})
	case "memory":
		Cacher = cache.NewMemoryCacher()
	default:
		log.Fatal().Msgf("unknown cache type: %s", config.Config.Cache.Type)
	}

	log.Debug().Str("type", config.Config.Cache.Type).Msg("init cache success")
	return
}
