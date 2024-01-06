package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var logger = log.New()

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(time.RFC3339)
	level := strings.ToUpper(entry.Level.String())
	message := entry.Message

	logString := fmt.Sprintf("[%s] %s: %s", timestamp, level, message)

	for key, value := range entry.Data {
		logString += fmt.Sprintf(" %s=%v", key, value)
	}

	logString += "\n"

	return []byte(logString), nil
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		logger.WithFields(log.Fields{
			"status_code": c.Writer.Status(),
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"ip":          c.ClientIP(),
		}).Info()
	}
}
