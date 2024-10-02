package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/zsais/go-gin-prometheus"
	"net/http"
	"time"
)

func StartGin(httpAddr string, logger log.Logger) error {
	r := gin.New()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	r.Use(gin.Logger())

	m := make(map[string]interface{})
	m["logger"] = logger
	r.Use(ConfigMiddleware(m))

	// 设置路由
	configRoutes(r)
	s := &http.Server{
		Addr:           httpAddr,
		Handler:        r,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		MaxHeaderBytes: 1 << 20,
	}
	level.Info(logger).Log("msg", "web_server_available_at", "httpAddr", httpAddr)

	err := s.ListenAndServe()
	return err
}
