package logger

import (
	"errors"
	"github.com/lookeme/short-url/internal/configuration"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

type Logger struct {
	Log *zap.Logger
}

var (
	LevelMap = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}
)

func CreateLogger(cfg *configuration.LoggerCfg) (*Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	level := zapcore.DebugLevel
	stdout := "stdout"
	if cfg.Level != "" {
		var ok bool
		level, ok = LevelMap[cfg.Level]
		if !ok {
			return nil, errors.New("unsupported logging level " + cfg.Level)
		}
	}
	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(level),
		OutputPaths:   []string{stdout},
		Encoding:      "json",
		EncoderConfig: encoderCfg,
	}
	logger := zap.Must(config.Build())
	return &Logger{
		Log: logger,
	}, nil
}

func (logger *Logger) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		logger.Log.Info("shorten path service ",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Duration("duration", duration),
		)
	}
	return http.HandlerFunc(fn)
}
