package middleware

import (
	"net/http"

	"github.com/kitabisa/perkakas/v2/log"
)

type HttpRequestLoggerMiddleware struct {
	logger *log.Logger
}

func NewHttpRequestLoggerMiddleware(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.SetRequest(r)
			next.ServeHTTP(w, r)
			logger.Print()
		})
	}
}
