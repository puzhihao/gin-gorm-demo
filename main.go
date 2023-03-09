package main

import (
	"github.com/puzhihao/gin-gorm-demo/app/host/impl"
	"github.com/puzhihao/gin-gorm-demo/conf"
	"github.com/puzhihao/gin-gorm-demo/routers"
)

func main() {

	conf.LoadConfigToml("./etc/restful-api.toml")
	impl.Service.Init()

	routers.InitRouter()

}
