package main

import (
	"fmt"
	"gin-hybrid/cmd"
	"gin-hybrid/conf"
	"gin-hybrid/rest"
	"gin-hybrid/router"
	"github.com/gin-gonic/gin"
)

func main() {
	srvConf := conf.MustNewServiceConfig[conf.Gateway]()
	restClient := rest.NewClient(srvConf)
	userService := restClient.MustAddServiceDependency("user")
	articleService := restClient.MustAddServiceDependency("article")
	fmt.Printf("%#v\n", srvConf)
	cmd.Entry(cmd.EntryConfig{Port: srvConf.SelfConf.Port}, func(engine *gin.Engine, api *gin.RouterGroup) {
		router.RegisterReverseProxy(userService, api.Group("/user"))
		router.RegisterReverseProxy(articleService, api.Group("/article"))
	})
}
