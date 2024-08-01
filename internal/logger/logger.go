package logger

import (
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lookeme/short-url/internal/configuration"
)

// Logger is a type that represents a logger instance.
type Logger struct {
	Log *zap.Logger
}

// LevelMap is a map that associates string keys with zapcore.Level values. It is used to map logging levels from string representations to their corresponding zapcore.Level constants
var (
	LevelMap = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}
)

// CreateLogger creates a new Logger instance based on the provided configuration.
// It takes a pointer to a LoggerCfg structure that specifies the logging configuration.
// The function returns a pointer to a Logger instance and an error if any occurred.
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

// Middleware is a method of the Logger struct that creates an HTTP middleware.
// It takes an http.Handler as input and returns an http.Handler.
// The returned handler performs logging for the incoming request and the corresponding response.
// It measures the request duration, captures the request URI, method, response status, and size.
// The captured information is then logged using the logger's Log.Info method.
// The capturing is done using a loggingResponseWriter, which wraps the original http.ResponseWriter.
// The captured response status and size are stored in a struct called responseData.
func (logger *Logger) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r)
		duration := time.Since(start)
		logger.Log.Info("shorten path service ",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Duration("duration", duration),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	}
	return http.HandlerFunc(fn)
}

// responseData is a type that represents the response data of an HTTP request.
//
// Fields:
// - status (int): the HTTP status code of the response.
// - size (int): the size of the response in bytes.
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter is a type that represents a response writer with logging capability.
// It embeds the original http.ResponseWriter and keeps track of the response data.
type loggingResponseWriter struct {
	http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
	responseData        *responseData
}

// Write is a method of the loggingResponseWriter struct.
// It
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

// WriteHeader is a method of the loggingResponseWriter struct that implements the http.ResponseWriter interface.
// It writes the provided status code to the response and captures it in the responseData struct.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}
