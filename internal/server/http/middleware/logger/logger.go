package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func New(log *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Collect data after request
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// Incoming request log
		log.Info().
			Str("client_ip", clientIP).
			Str("path", path).
			Str("method", method).
			Msg("incoming request")

		// Process request
		c.Next()

		latency := time.Since(start)
		// Outgoing request log
		log.Debug().
			Dur("latency", latency).
			Int("status", status).
			Str("client_ip", clientIP).
			Str("path", path).
			Str("method", method).
			Msg("request completed")

	}
}
