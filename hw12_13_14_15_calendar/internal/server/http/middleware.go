package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Запоминаем время начала обработки запроса
		start := time.Now()

		// Обработать запрос
		next.ServeHTTP(w, r)

		// Получаем информацию о статусе ответа
		statusCode := w.Header().Get("Status")

		if statusCode == "" {
			statusCode = "200"
		}

		latency := time.Since(start)
		userAgent := r.UserAgent()

		logger.Info(fmt.Sprintf("%s [%s] %s %s %s %s %d \"%s\"\n",
			r.RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			statusCode,
			latency.Milliseconds(),
			userAgent,
		))
	})
}
