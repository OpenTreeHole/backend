package config

import (
	. "github.com/opentreehole/backend/common"
	"github.com/spf13/viper"
)

var defaultConfig = map[string]string{
	EnvMode:      "dev",
	EnvPort:      "8000",
	EnvDBType:    "sqlite",
	EnvCacheType: "memory",
}

func init() {
	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}
