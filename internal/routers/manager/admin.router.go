package manager

import "github.com/gin-gonic/gin"

type AdminRouter struct{}

func (ar *AdminRouter) InitAdminRouter(Router *gin.RouterGroup) {
	// public router
	adminRouterPublic := Router.Group("/admin")
	{
		adminRouterPublic.POST("/login")
		adminRouterPublic.GET("/otp")
	}

	// private router
	// adminRouterPrivate := Router.Group("/admin/user")
	// // userRouterPrivate.Use(Limiter()) // Add any necessary middleware here, e.g., authentication
	// // userRouterPrivate.Use(Authen()) // Add any necessary middleware here, e.g., authentication
	// // userRouterPrivate.Use(Permission()) // Add any necessary middleware here, e.g., authentication
	// {
	// 	adminRouterPrivate.POST("/active_user")
	// }

}
