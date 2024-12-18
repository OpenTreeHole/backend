package config

import (
	"github.com/spf13/viper"
)

const (
	EnvHostName = "HOST_NAME"
)

var defaultConfig = map[string]string{
	EnvHostName: "localhost:8000",
}

func init() {
	viper.AutomaticEnv()
	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}
