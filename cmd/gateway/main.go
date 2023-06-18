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
	srvConf, err := conf.NewServiceConfig[conf.Gateway]("gateway")
	if err != nil {
		panic(err)
	}
	conf.GatewayServiceConfig = srvConf
	restClient := rest.NewClient(srvConf.Etclient)
	_, err = restClient.AddService("user")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", srvConf)
	cmd.Entry(cmd.EntryConfig{Port: srvConf.SelfConf.Port}, func(engine *gin.Engine, api *gin.RouterGroup) {
		router.RegisterAPIRouters(router.GetUserAPIRouters(), api.Group("/user"))
	})
}
