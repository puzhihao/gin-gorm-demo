package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/puzhihao/gin-gorm-demo/app/host/http"
)

func InitRouter() {
	r := gin.Default()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/host", http.GetHost)
	r.POST("/host", http.CreateHost)

	r.Run(":8080")
}
