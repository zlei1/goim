package http

import (
	"fmt"
	"net/http/httputil"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func loggerHandler(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	method := c.Request.Method

	c.Next()

	end := time.Now()
	latency := end.Sub(start)
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()
	if raw != "" {
		path = path + "?" + raw
	}
	log.WithFields(log.Fields{
		"METHOD": method,
		"PATH": path,
		"CODE": statusCode,
		"IP": clientIP,
		"TIME": latency/time.Millisecond,
	}).Info("HTTP REQUEST")
}

func recoverHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			httpRequest, _ := httputil.DumpRequest(c.Request, false)
			pnc := fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s", time.Now().Format("2006-01-02 15:04:05"), string(httpRequest), err, buf)
			log.Error(pnc)
			c.AbortWithStatus(500)
		}
	}()
	c.Next()
}