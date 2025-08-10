package routers

import "github.com/gin-gonic/gin"

func Init(gin *gin.Engine) *gin.Engine {
	freeRoutes := gin.Group("")

	AuthRoutes(freeRoutes)
	return gin
}