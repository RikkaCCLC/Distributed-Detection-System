package web

import "github.com/gin-gonic/gin"

func ConfigMiddleware(m map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		for k, v := range m {
			c.Set(k, v)
			c.Next()
		}
	}
}
