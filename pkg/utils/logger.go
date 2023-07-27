package utils

import (
	"time"

	"github.com/rs/zerolog"
)

func init() {
	// compatible with zap and old spec
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339Nano
}
