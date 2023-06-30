package common

import (
	"github.com/rs/zerolog"
	"time"
)

func init() {
	// compatible with zap and old spec
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339Nano
}
