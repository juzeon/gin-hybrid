package main

import (
	"gin-hybrid/cmd"
	"gin-hybrid/conf"
	"gin-hybrid/rest"
	"gin-hybrid/router"
	"github.com/gin-gonic/gin"
)

func main() {
	srvConf := conf.MustNewServiceConfig[conf.Article]()
	restClient := rest.NewClient(srvConf)
	userService := restClient.MustAddServiceDependency("user")
	cmd.Entry(cmd.EntryConfig{Port: srvConf.SelfConf.Port}, func(engine *gin.Engine, api *gin.RouterGroup) {
		router.RegisterAPIRouters(getAPIRouters(userService), api, srvConf)
	})
}

func getAPIRouters(userService *rest.Service) []router.APIRouter {
	srv := NewArticleService(userService)
	apiRouters := []router.APIRouter{
		{
			Method:   "post",
			Path:     "/articles",
			Handlers: router.AssembleHandlers(srv.PostArticle),
		},
		{
			Method:   "get",
			Path:     "/articles",
			Handlers: router.AssembleHandlers(srv.ListArticle),
		},
	}
	return apiRouters
}
