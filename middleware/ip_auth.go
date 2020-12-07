package middleware

import (
	"gin-cladder/conf/elite/control"
	"github.com/gin-gonic/gin"
)

func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isMatched := false
		for _,host := range control.GetStringSliceConf(control.ConfEnv,"http.allow_host") {
			if c.ClientIP() == "::1" ||c.ClientIP() ==  host {
				isMatched = true
			}
		}
		if !isMatched {
			c.Abort()
			return
		}
		c.Next()
	}
}