package rest_api

import (
	"github.com/VSKrivoshein/FBS-test/internal/app/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
)

// JSONLogMiddleware logs a gin HTTP requests in JSON format, with some additional custom key/values
func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := utils.GetDurationInMilliseconds(start)
		entry := log.WithFields(log.Fields{
			"method":   c.Request.Method,
			"path":     c.Request.RequestURI,
			"status":   c.Writer.Status(),
			"duration": duration,
		})

		if c.Writer.Status() >= 500 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info(c.Errors.String())
		}
	}
}
