package log

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const KEY = "zapLogger"

type Logger struct {
	*zap.Logger
}

func NewLogger(conf *viper.Viper) (*Logger, error) {
	zapLogger, err := initZap(conf)
	if err != nil {
		return nil, err
	}
	return &Logger{Logger: zapLogger}, nil
}

func initZap(config *viper.Viper) (*zap.Logger, error) {
	var (
		atomicLevel zapcore.Level
		development = config.GetString("mode") != "production"
	)
	if config.GetString("mode") != "production" {
		atomicLevel = zapcore.DebugLevel
	} else {
		atomicLevel = zapcore.InfoLevel
	}
	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(atomicLevel),
		Development: development,
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
