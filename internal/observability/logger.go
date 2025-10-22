package observability

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func InitLogger(service string, level string) (*zap.Logger, error) {
	lvl := zapcore.InfoLevel
	if err := lvl.UnmarshalText([]byte(strings.ToLower(level))); err != nil {
		lvl = zapcore.InfoLevel
	}
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(lvl),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg", LevelKey: "level", TimeKey: "ts",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			CallerKey:   "caller", EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	l, err := cfg.Build(zap.AddCaller(), zap.Fields(zap.String("service", service)))
	if err != nil {
		return nil, err
	}
	log = l
	return log, nil
}

func L() *zap.Logger {
	if log == nil {
		_, _ = InitLogger(os.Getenv("SERVICE_NAME"), os.Getenv("LOG_LEVEL"))
	}
	return log
}
