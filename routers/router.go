package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/puzhihao/gin-gorm-demo/app/host/http"
)

func InitRouter() {
	r := gin.Default()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/hosts", http.GetHost)
	r.POST("/hosts", http.CreateHost)
	r.GET("/hosts/:id", http.DescribeHost)
	r.PATCH("/hosts/:id", http.PatchHost)
	r.PUT("/hosts/:id", http.UpdateHost)
	r.DELETE("/hosts/:id", http.DeleteHost)

	r.Run(":8080")
}
