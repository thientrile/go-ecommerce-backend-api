package manager
import "github.com/gin-gonic/gin"

type UserRouter struct{}

func (us *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// // public router
	// userRouterPublic := Router.Group("/admin/user")
	// {
	// 	userRouterPublic.POST("/register")
	// 	userRouterPublic.GET("/otp")
	// }

	// private router
	userRouterPrivate := Router.Group("/admin/user")
	// userRouterPrivate.Use(Limiter()) // Add any necessary middleware here, e.g., authentication
	// userRouterPrivate.Use(Authen()) // Add any necessary middleware here, e.g., authentication
	// userRouterPrivate.Use(Permission()) // Add any necessary middleware here, e.g., authentication
	{
		userRouterPrivate.POST("/active_user")
	}

}