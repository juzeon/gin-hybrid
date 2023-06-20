package main

import (
	"gin-hybrid/cmd"
	"gin-hybrid/conf"
	"gin-hybrid/middleware"
	"gin-hybrid/router"
	"gin-hybrid/service"
	"github.com/gin-gonic/gin"
)

func main() {
	srvConf := conf.MustNewServiceConfig[conf.User]()
	cmd.Entry(cmd.EntryConfig{Port: srvConf.SelfConf.Port},
		func(engine *gin.Engine, api *gin.RouterGroup) {
			router.RegisterAPIRouters(getAPIRouters(), api, srvConf)
		})
}

func getAPIRouters() []router.APIRouter {
	srv := service.NewUserService()
	routers := []router.APIRouter{
		{
			Method:   "post",
			Path:     "/login",
			Handlers: router.AssembleHandlers(srv.Login),
		},
		{
			Method:   "get",
			Path:     "/me",
			Handlers: router.AssembleHandlers(middleware.Auth, srv.Me),
		},
		{
			Method:   "get",
			Path:     "/example",
			Handlers: router.AssembleHandlers(srv.ExampleGet),
			RPCOnly:  true,
		},
		{
			Method:   "post",
			Path:     "/example",
			Handlers: router.AssembleHandlers(srv.ExamplePost),
			RPCOnly:  true,
		},
	}
	return routers
}
