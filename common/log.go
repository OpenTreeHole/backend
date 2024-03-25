package common

import (
	"log/slog"
	"os"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type mLogger struct {
	*slog.Logger
}

func (l mLogger) Printf(message string, args ...any) {
	l.Info(message, args...)
}
