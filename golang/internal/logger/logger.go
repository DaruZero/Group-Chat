package logger

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new logger
func New(logLevelEnv string) *zap.SugaredLogger {
	logLevel := os.Getenv(logLevelEnv)
	if logLevel == "" {
		logLevel = "info"
	}

	encoderConfig := ecszap.NewDefaultEncoderConfig()
	var core zapcore.Core

	switch logLevel {
	case "debug":
		core = ecszap.NewCore(
			encoderConfig, os.Stdout, zap.DebugLevel)
	case "info":
		core = ecszap.NewCore(
			encoderConfig, os.Stdout, zap.InfoLevel)
	case "warn":
		core = ecszap.NewCore(
			encoderConfig, os.Stdout, zap.WarnLevel)
	case "error":
		core = ecszap.NewCore(
			encoderConfig, os.Stdout, zap.ErrorLevel)
	case "fatal":
		core = ecszap.NewCore(
			encoderConfig, os.Stdout, zap.FatalLevel)
	default:
		core = ecszap.NewCore(
			encoderConfig, os.Stdout, zap.InfoLevel)
	}

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)

	return logger.Sugar()
}
