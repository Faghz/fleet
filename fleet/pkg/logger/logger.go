package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(env string, level string) (*zap.Logger, error) {
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	if env == "production" {
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(logLevel)
		logger, err := config.Build(zap.AddStacktrace(zap.PanicLevel))

		if err != nil {
			return nil, fmt.Errorf("failed to create production logger: %w", err)
		}

		return logger, nil
	}

	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(logLevel)
	logger, err := config.Build(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		return nil, fmt.Errorf("failed to create development logger: %w", err)
	}

	return logger, nil
}
