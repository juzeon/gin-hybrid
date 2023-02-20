package main

import (
	"fmt"
	"gin-hybrid/router"
	"gin-hybrid/service"
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	engine.Use(gin.Logger(), nice.Recovery(router.RecoveryFunc))
	service.Setup()
	router.Setup(engine)
	err := engine.Run(fmt.Sprintf(":%v", 7070))
	if err != nil {
		panic(err)
	}
}
