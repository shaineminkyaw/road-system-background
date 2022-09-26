package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaineminkyaw/road-system-background/config"
	"github.com/sirupsen/logrus"
)

func Logging() gin.HandlerFunc {
	logger := config.LogConf()

	return func(c *gin.Context) {
		//Start time
		startTime := time.Now()

		//Process request
		c.Next()

		//End time
		// endTime := time.Now()

		//Execution time
		latencyTime := time.Since(startTime).Nanoseconds()

		//Request method
		reqMethod := c.Request.Method

		//Request routing
		reqUri := c.Request.RequestURI

		// status code
		statusCode := c.Writer.Status()

		// request IP
		clientIP := c.ClientIP()

		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()
	}
}
