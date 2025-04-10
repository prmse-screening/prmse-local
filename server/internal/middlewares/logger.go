package middlewares

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		var logBuffer bytes.Buffer
		logBuffer.WriteString(fmt.Sprintf("%v | %s | %s | %s",
			duration,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		))

		errMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if errMsg != "" {
			logBuffer.WriteString(fmt.Sprintf(" | %s", errMsg))
		}

		log.Info(logBuffer.String())
	}
}
