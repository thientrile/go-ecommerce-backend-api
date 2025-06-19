package initialize

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-ecommerce-backend-api.com/cmd/swag/docs"
	"go-ecommerce-backend-api.com/global"
)

func InitSwagger(r *gin.Engine) *gin.Engine {

	docs.SwaggerInfo.Title = "API Documentation Ecommerce Backend"
	docs.SwaggerInfo.Description = "This is a sample server celler server."
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.Host = global.Config.Server.Domain
	docs.SwaggerInfo.BasePath = "/v1/2025"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
