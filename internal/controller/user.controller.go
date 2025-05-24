package controller

import (

	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/internal/service"
	"go-ecommerce-backend-api.com/pkg/response"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func (uc *UserController) GetUserById(c *gin.Context){
	// Get the user ID from the URL parameter
	// userId := c.Param("id")

	// // Simulate fetching user data from a database or service
	// user := map[string]string{
	// 	"id":   userId,
	// 	"name": "John Doe",
	// }

	// // Return the user data as JSON
	// c.JSON(200, gin.H{
	// 	"user": user,
	// })


	// c.JSON(http.StatusOK, response.ResponstData{
	// 	Code:2001,
	// 	Message: "Success",
	// 	Data: []string{"hello", "world"},

	// })
	response.SuccessResponse(c, 20001, uc.userService.GetInfoUser())
}