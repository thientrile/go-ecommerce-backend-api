package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/internal/service"
	"go-ecommerce-backend-api.com/internal/vo"
	"go-ecommerce-backend-api.com/pkg/response"
)

// interface UserControllerInterface

type UserController struct {
	userService service.IUserService
}

func NewUserController(
	userService service.IUserService,
) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) Regisger(c *gin.Context) {
	var params vo.UserRegistratorRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamInvalid, err.Error())
		return
	}
	fmt.Println("Email params:", params.Email)
	result := uc.userService.Regisger(params.Email, params.Purpose)
	response.SuccessResponse(c, result, nil)
}
