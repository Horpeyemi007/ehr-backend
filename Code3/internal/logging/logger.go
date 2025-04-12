package logging

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

// This sets up the global logger
func InitializeLogger() {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.DisableStacktrace = true
	logger, _ := zapConfig.Build()
	Logger = logger.Sugar()
}

// This ensures that logs are flushed before the app exits
func CleanupLogger() {
	if Logger != nil {
		Logger.Sync()
	}
}
