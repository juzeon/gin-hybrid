package main

import (
	"gin-hybrid/cmd"
	"gin-hybrid/conf"
	"gin-hybrid/router"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	conf.GatewayConf.Load()
	log.Printf("%#v", conf.InitConf)
	log.Printf("%#v", conf.ParentConf)
	log.Printf("%#v", conf.GatewayConf)
	cmd.Entry(cmd.EntryConfig{Port: conf.GatewayConf.Port}, func(engine *gin.Engine, api *gin.RouterGroup) {
		router.RegisterAPIRouters(router.GetUserAPIRouters(), api.Group("/user"))
	})
}
