package cmd

import (
	"fmt"
	"gin-hybrid/router"
	"gin-hybrid/service"
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
)

func Entry(registerFunc func(engine *gin.Engine)) {
	engine := gin.New()
	engine.Use(gin.Logger(), nice.Recovery(router.RecoveryFunc))
	service.Setup()
	router.Setup(engine)
	registerFunc(engine)
	err := engine.Run(fmt.Sprintf(":%v", 7070))
	if err != nil {
		panic(err)
	}
}
