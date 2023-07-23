package zapctx

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// This is required to make the logger available in the ctx or fetch it from the ctx
var loggerKey struct{}

// NewLogger returns a new zap logger. This logger is configured to log at InfoLevel and above
// We are using this logger instead of zap.NewProduction because by default it's configured to sample the logs
// Additionally, we this log will not fail to be created which simplifies what things could go wrong
func NewLogger() *zap.Logger {
	return NewLoggerWithConfig(DefaultLoggerConfig())
}

// NewLoggerWithConfig returns a new zap logger. This logger would use the config provided by the caller.
// You're only able to configure to set the JSON EncoderWithLogs and the level. The output is defaulted to stdout.
func NewLoggerWithConfig(config zap.Config) *zap.Logger {
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			os.Stdout,
			config.Level,
		),
	)
}

func DefaultLoggerConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	return cfg
}

// WithLogger attaches a zap logger to a given context, returning the new updated context
func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// Logger extracts a zap logger from a context, providing a default logger if one is not present
// It's recommended that you do not change the logger config too much from the default because when you use this method
// we will not respect that configuration
func Logger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		panic("nil context passed to Logger")
	}

	if logger, _ := ctx.Value(loggerKey).(*zap.Logger); logger != nil {
		return logger
	}

	return NewLogger()
}
