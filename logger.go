package common

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LoggerOld = zerolog.New(os.Stdout).With().Timestamp().Logger()

func init() {
	// compatible with zap and old spec
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

const KEY = "zapLogger"

type Logger struct {
	*zap.Logger
}

func NewLogger() (*Logger, func()) {
	// Info Level
	logger, err := initZap()
	if err != nil {
		panic(err)
	}
	return &Logger{Logger: logger}, func() {
		_ = logger.Sync()
	}
}

func initZap() (*zap.Logger, error) {
	var atomicLevel zapcore.Level
	// Info Level, production env
	atomicLevel = zapcore.InfoLevel

	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(atomicLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return logConfig.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
}

// NewContext Adds a field to the specified context
func (l *Logger) NewContext(c *fiber.Ctx, fields ...zapcore.Field) {
	c.Locals(KEY, &Logger{l.WithContext(c).With(fields...)})
}

// WithContext Returns a zap instance from the specified context
func (l *Logger) WithContext(c *fiber.Ctx) *Logger {
	if c == nil {
		return l
	}
	ctxLogger, ok := c.Locals(KEY).(*Logger)
	if ok {
		return ctxLogger
	}
	return l
}
