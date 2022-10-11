package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/woodlsy/woodGin/log"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				//httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					log.Logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						//zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					log.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						//zap.String("request", string(httpRequest)),

					)
					log.Logger.Error(zap.String("stack", string(debug.Stack())))
				} else {
					log.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						//zap.String("request", string(httpRequest)),
					)
				}
				c.JSON(http.StatusOK, "系统错误，请联系管理员")
			}
		}()
		c.Next()

	}
}
