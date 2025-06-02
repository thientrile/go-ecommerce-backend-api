package controller

import (
	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/internal/service"
	"go-ecommerce-backend-api.com/pkg/response"
)

// interface UserControllerInterface


type UserController struct{
	userService service.IUserService
}

func NewUserController(
	userService service.IUserService,
) *UserController {
	return &UserController{
		userService: userService,
	}
}


func (uc *UserController) Regisger(c *gin.Context){
	result := uc.userService.Regisger("", "")
	response.SuccessResponse(c,result,nil)
}
