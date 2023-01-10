package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
)

func logError(err error) {
	if err == nil {
		return
	}

	var e *errors.Error
	if errors.As(err, &e) {
		logger.Errorf("HTTP Error: %v", e.ErrorStack())
	} else {
		logger.Errorf("HTTP Error: %v", err)
	}
}

func httpLogger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	c.Next()

	timeStamp := time.Now()
	millis := timeStamp.Sub(start).Milliseconds()

	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	respSize := c.Writer.Size()
	if raw != "" {
		path = path + "?" + raw
	}

	logger.Infof("%s %d %s (resp: %d) (took: %dms) (remote: %s)", method, statusCode, path, respSize, millis, clientIP)

	for _, err := range c.Errors {
		logError(err)
	}

	// c.Errors
}
