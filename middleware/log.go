package middleware

import (
	"net/http"

	"github.com/kitabisa/perkakas/v2/log"
)

type HttpRequestLoggerMiddleware struct {
	logger *log.Logger
}

func NewHttpRequestLoggerMiddleware(logger *log.Logger) *HttpRequestLoggerMiddleware {
	return &HttpRequestLoggerMiddleware{
		logger: logger,
	}
}

func (l *HttpRequestLoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	l.logger.SetRequest(r)
	next(w, r)
	l.logger.Print()
}
