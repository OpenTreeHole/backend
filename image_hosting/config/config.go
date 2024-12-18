package config

import (
	"context"
	"github.com/opentreehole/backend/common"
	"github.com/spf13/viper"
	"log/slog"
)

var (
	DbUrl    string
	HostName string
)

func init() {
	DbUrl = viper.GetString(common.EnvDBUrl)
	if DbUrl == "" {
		slog.LogAttrs(context.Background(), slog.LevelError, "", slog.String("err", "DB_URL is empty"))
	}
	HostName = viper.GetString("HOST_NAME")
	if HostName == "" {
		// default value
		HostName = "localhost:8000"
	}
}
