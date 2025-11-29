package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a configured Zap logger.
// Returns the logger and a cleanup function to flush buffers
func NewLogger(env string) (*zap.Logger, func(), error) {
	var config zap.Config

	if env == "local" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	logger, err := config.Build()
	if err != nil {
		return nil, nil, err
	}

	// cleanup function ensures logs are written before the app exits
	cleanup := func() {
		_ = logger.Sync()
	}

	return logger, cleanup, nil
}
