package main

import (
	"gin-hybrid/cmd"
	"gin-hybrid/router"
	"github.com/gin-gonic/gin"
)

func main() {
	cmd.Entry(func(engine *gin.Engine) {
		api := engine.Group("/api")
		router.RegisterAPIRouters(router.GetUserAPIRouters(), api.Group("/user"))
	})
}
