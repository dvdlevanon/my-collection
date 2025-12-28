package server

import (
	"context"
	"my-collection/server/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

type Handler interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

func shrinkPath(method, path string, maxLen int) string {
	if path == "" || path == " " {
		return method
	}

	// Strip common prefix
	path = strings.TrimPrefix(path, "/api/v1")
	path = strings.TrimPrefix(path, "/api")
	if path == "" {
		path = "/" // Handle edge case where path was exactly "/api"
	}

	full := method + " " + path
	if len(full) <= maxLen {
		return full // Don't shrink if not needed
	}

	// Split path into parts
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return method
	}

	// Try progressive shrinking - start from the left
	current := path
	for i := 0; i < len(parts); i++ {
		// Build path with first i parts as acronyms, rest as full
		var pathParts []string
		for j := 0; j < len(parts); j++ {
			if j <= i && len(parts[j]) > 0 {
				pathParts = append(pathParts, strings.ToUpper(string(parts[j][0])))
			} else {
				pathParts = append(pathParts, parts[j])
			}
		}
		current = "/" + strings.Join(pathParts, "/")
		result := method + " " + current

		if len(result) <= maxLen {
			return result
		}
	}

	// If still too long, full acronym
	acronym := ""
	for _, part := range parts {
		if len(part) > 0 {
			acronym += strings.ToUpper(string(part[0]))
		}
	}
	result := method + " " + acronym

	// Last resort: truncate
	if len(result) > maxLen {
		return result[:maxLen]
	}

	return result
}
func ContextWithSubject(c *gin.Context) context.Context {
	ctx := c.Request.Context()

	// If already has subject, return as-is (don't increment counter again)
	if utils.HasSubject(ctx) {
		return ctx
	}

	// Initialize for the first time
	endpoint := c.Request.Method + " " + c.FullPath()
	if endpoint == " " {
		endpoint = c.Request.Method + " " + c.Request.URL.Path
	}

	endpoint = shrinkPath(c.Request.Method, c.FullPath(), 20)
	return utils.ContextWithSubject(ctx, endpoint)
}
func HandleError(ctx context.Context, c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	httpError := http.StatusInternalServerError
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpError = http.StatusNotFound
	}

	utils.LogError(ctx, logger, "Server error", err)
	c.AbortWithError(httpError, err)
	return true
}

func logError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	var e *errors.Error
	if errors.As(err, &e) {
		logger.Errorf(ctx, "HTTP Error: %v", e.ErrorStack())
	} else {
		logger.Errorf(ctx, "HTTP Error: %v", err)
	}
}

func contextInitializer(c *gin.Context) {
	ctx := ContextWithSubject(c)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func httpLogger(c *gin.Context) {
	ctx := ContextWithSubject(c)
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

	logger.Infof(ctx, "%s %d %s (resp: %d) (took: %dms) (remote: %s)", method, statusCode, path, respSize, millis, clientIP)

	for _, err := range c.Errors {
		logError(ctx, err)
	}
}
