package utils

import (
	"testing"

	"github.com/rs/zerolog/log"
)

func TestLog(t *testing.T) {
	log.Print("hello world")

	log.Info().Msg("hello world")
}
