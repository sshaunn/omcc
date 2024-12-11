package middleware

import (
	"github.com/gin-gonic/gin"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"time"
)

func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		elapseTime := time.Since(start)
		log.Info("Received http request",
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.String("query", raw),
			logger.Int("status", c.Writer.Status()),
			logger.Duration("elapseTime", elapseTime),
			logger.String("ip", c.ClientIP()),
		)
	}
}
