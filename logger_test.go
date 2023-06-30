package common

import (
	"github.com/rs/zerolog/log"
	"testing"
)

func TestLog(t *testing.T) {
	log.Print("hello world")

	log.Info().Msg("hello world")
}
