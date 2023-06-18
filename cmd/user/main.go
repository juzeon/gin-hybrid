package main

import (
	"gin-hybrid/cmd"
	"gin-hybrid/conf"
	"gin-hybrid/router"
	"github.com/gin-gonic/gin"
)

func main() {
	srvConf, err := conf.NewServiceConfig[conf.User]("user")
	if err != nil {
		panic(err)
	}
	conf.UserServiceConfig = srvConf
	cmd.Entry(cmd.EntryConfig{Port: srvConf.SelfConf.Port}, func(engine *gin.Engine, api *gin.RouterGroup) {
		router.RegisterAPIRouters(router.GetUserAPIRouters(), api.Group("/user"))
	})
}
