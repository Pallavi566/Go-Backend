package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new logger instance
func NewLogger(serviceName, environment string) (*zap.Logger, error) {
	// Configure the encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Set the log level based on environment
	var level zapcore.Level
	switch environment {
	case "production":
		level = zap.InfoLevel
	default:
		level = zap.DebugLevel
	}

	// Create a config with the encoder config and level
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      environment != "production",
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"service":   serviceName,
			"env":       environment,
		},
	}

	// Add caller skip for production to get the correct caller info
	if environment == "production" {
		return config.Build(zap.AddCallerSkip(1))
	}

	return config.Build()
}

// NewTestLogger creates a logger suitable for testing
func NewTestLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"/dev/null"}
	logger, _ := config.Build()
	return logger
}
