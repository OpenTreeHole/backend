package config

import "github.com/caarlos0/env/v9"

var Config struct {
	DbURL string `env:"DB_URL"`
}

func init() {
	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}
}
