package cmd

import (
	"fmt"
	"gin-hybrid/router"
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

type EntryConfig struct {
	Port int
}

func Entry(entryConfig EntryConfig, registerFunc func(engine *gin.Engine, api *gin.RouterGroup)) {
	rand.Seed(time.Now().Unix())
	engine := gin.New()
	engine.Use(gin.Logger(), nice.Recovery(router.RecoveryFunc))
	api := engine.Group("/api")
	registerFunc(engine, api)
	router.Setup(engine, false)
	err := engine.Run(fmt.Sprintf(":%v", entryConfig.Port))
	if err != nil {
		panic(err)
	}
}