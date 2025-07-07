package user

import (
	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/internal/controller/ticket"
)

type TicketRouter struct{}

func (tr *TicketRouter) InitTicketRouter(Router *gin.RouterGroup) {
	// public router
	ticketRouterPublic := Router.Group("/ticket")
	{
		ticketRouterPublic.GET("/search")                                        // Search tickets
		ticketRouterPublic.GET("/item/:id", ticket.TicketItem.GetTicketItemById) // Get ticket details by ID
	}


}
