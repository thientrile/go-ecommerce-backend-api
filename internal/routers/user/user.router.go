package user

import "github.com/gin-gonic/gin"

type UserRouter struct{}

func (us *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// public router
	userRouterPublic := Router.Group("/user")
	{
		userRouterPublic.POST("/register") // Register a new user -> Yes ->No
		userRouterPublic.GET("/otp")
	}

	// private router
	userRouterPrivate := Router.Group("/user")
	// userRouterPrivate.Use(Limiter()) // Add any necessary middleware here, e.g., authentication
	// userRouterPrivate.Use(Authen()) // Add any necessary middleware here, e.g., authentication
	// userRouterPrivate.Use(Permission()) // Add any necessary middleware here, e.g., authentication
	{
		userRouterPrivate.GET("/get_info")
	}

}
