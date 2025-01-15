package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}

	Logger struct {
		zap   *zap.Logger
		Sugar *zap.SugaredLogger
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.responseData.status = status
}

func NewLogger() (*Logger, error) {
	zap, err := zap.NewDevelopment()
	return &Logger{zap: zap, Sugar: zap.Sugar()}, err
}

func (l *Logger) Close() {
	l.zap.Sync()
}

func (l *Logger) WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

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

		var logMsg []interface{}

		switch r.Method {
		case http.MethodGet:
			logMsg = append(logMsg, "URI", r.RequestURI, "  METHOD:", r.Method, "  DURATION:", duration)
		case http.MethodPost:
			logMsg = append(logMsg, "URI", r.RequestURI, "  METHOD:", r.Method,
				"  DURATION:", duration, "  STATUS", responseData.status, "  SIZE", responseData.size)
		}

		l.Sugar.Infoln(logMsg...)
	})
}
