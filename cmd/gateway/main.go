package main

import (
	"fmt"
	"gin-hybrid/cmd"
	"gin-hybrid/conf"
	"gin-hybrid/router"
	"github.com/gin-gonic/gin"
)

func main() {
	conf.GatewayConf.Load()
	fmt.Println(conf.InitConf)
	fmt.Println(conf.ParentConf)
	fmt.Println(conf.GatewayConf)
	cmd.Entry(cmd.EntryConfig{Port: conf.GatewayConf.Port}, func(engine *gin.Engine, api *gin.RouterGroup) {
		router.RegisterAPIRouters(router.GetUserAPIRouters(), api.Group("/user"))
	})
}
