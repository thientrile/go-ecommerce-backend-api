package initialize

import (
	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/internal/middlewares"
	"go-ecommerce-backend-api.com/internal/routers"
)

func InitRouter() *gin.Engine {

	var r *gin.Engine
	s := global.Config.Server
	if s.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
		// Use custom logger middleware instead of default gin logger
		// r.Use(gin.Recovery())
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		// Use custom logger middleware
		// r.Use(gin.Recovery())
	}
	// middlewares
	// r.Use() //cross origin resource sharing

	// Initialize rate limiter once

	r.Use(middlewares.NewRateLimiter().DynamicRateLimiter()) // limiter global.Config.Limiter
	r.Use(middlewares.LoggerMiddleware())                    // logger - đã sử dụng custom logger middleware ở trên
	managerRouter := routers.RouterGroupApp.Manager
	userRouter := routers.RouterGroupApp.User

	MainGroup := r.Group("/v1/2025")
	{
		MainGroup.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "pong",
			})
		}) // tracking monitoring
	}
	{
		// user router
		userRouter.InitUserRouter(MainGroup)
		userRouter.InitProductRouter(MainGroup)

	}
	{
		// manager router
		managerRouter.InitUserRouter(MainGroup)
		managerRouter.InitAdminRouter(MainGroup)
	}

	// Initialize Swagger documentation

	return r
}
