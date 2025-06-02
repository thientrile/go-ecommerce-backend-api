package user

import (
	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/internal/controller"
	"go-ecommerce-backend-api.com/internal/repo"
	"go-ecommerce-backend-api.com/internal/service"
)

type UserRouter struct{}

func (pr *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// public router
	// this non-dependency
	ur := repo.NewUserRepository()
	us := service.NewUserService(ur)
	useHandlerNoDependency := controller.NewUserController(us)

	// Wire go


	
	userRouterPublic := Router.Group("/user")
	{
		userRouterPublic.POST("/register",useHandlerNoDependency.Regisger) // Register a new user -> Yes ->No
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
