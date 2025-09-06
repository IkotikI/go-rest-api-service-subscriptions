package handler

import (
	_ "github.com/ikotiki/go-rest-api-service-subscriptions/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type SwaggerController struct{}

func NewSwaggerController(g *gin.RouterGroup) *SwaggerController {
	a := &SwaggerController{}
	a.registerRoutes(g)
	return a
}

func (a *SwaggerController) registerRoutes(g *gin.RouterGroup) {
	// g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	g.GET("/swagger/*any", func(c *gin.Context) {
		// Disable caching
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

}
