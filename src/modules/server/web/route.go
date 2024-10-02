package web

import (
	"github.com/gin-gonic/gin"
	"time"
)

func configRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})
		api.GET("/now-ts", GetNowTs)
		api.POST("/node-path", NodePathAdd)
		api.GET("/node-path", NodePathQuery)
		api.POST("/resource-mount", ResourceMount)
		api.DELETE("/resource-unmount", ResourceUnMount)
		api.POST("/resource-query", ResourceQuery)
		api.POST("/resource-distribution", GetLabelDistribution)
		api.GET("/resource-group", ResourceGroup)
		api.POST("/log-job", LogJobAdd)
		api.GET("/log-job", LogJobGets)

		api.POST("/task", TaskAdd)
		api.GET("/task", TaskGets)
		api.POST("/kill-task", TaskKill)

	}
}

func GetNowTs(c *gin.Context) {
	c.String(200, time.Now().Format("2006-01-02 15:04:05"))
}
