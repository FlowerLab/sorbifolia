//go:build gin

package rogu

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.x2ox.com/sorbifolia/coarsetime"
)

// Gin returns a gin.HandlerFunc middleware that logs requests.
func Gin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			start = coarsetime.FloorTime()
			path  = c.Request.URL.Path
			query = c.Request.URL.RawQuery
		)

		c.Next()
		end := coarsetime.FloorTime()

		fields := []zapcore.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", end.Sub(start)),
		}
		if e := c.Errors.Errors(); len(e) != 0 {
			fields = append(fields, zap.Strings("errors", e))
			zap.L().Error(path, fields...)
			return
		}
		zap.L().Info(path, fields...)
	}
}

func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				zap.L().Error("GinRecovery",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(stack())),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
