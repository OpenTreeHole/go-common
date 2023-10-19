package common

import (
	"testing"

	"github.com/rs/zerolog/log"
)

func TestLog(t *testing.T) {
	log.Print("hello world")

	log.Info().Msg("hello world")

	LoggerOld.Info().Msg("hello world")
}
