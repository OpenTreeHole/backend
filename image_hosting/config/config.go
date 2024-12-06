package config

import "github.com/caarlos0/env/v9"

var Config struct {
	DbURL    string `env:"DB_URL"`
	HostName string `env:"HOST_NAME" envDefault:"localhost:8000"`
}

func init() {
	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}
}
