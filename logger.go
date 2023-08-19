package common

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func init() {
	// compatible with zap and old spec
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339Nano
}
