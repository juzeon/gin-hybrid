package app

import "github.com/gin-gonic/gin"

func HandlerFunc(handler func(aw *Wrapper) Result) func(c *gin.Context) {
	return func(c *gin.Context) {
		aw := NewWrapper(c)
		res := handler(aw)
		c.JSON(200, res)
	}
}
