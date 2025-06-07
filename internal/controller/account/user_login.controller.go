package account

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/internal/model"
	"go-ecommerce-backend-api.com/internal/service"
	"go-ecommerce-backend-api.com/pkg/response"
	"go.uber.org/zap"
)

// manager controller login user
var Login = new(cUserLogin)

type cUserLogin struct {
}

func (c *cUserLogin) Login(ctx *gin.Context) {
	// Implement login logic
	err := service.UserLogin().Login(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidOTP, err.Error())
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)

}

// User Register documentation
//
//	@Summary      Register user with OTP
//	@Description  Register a new user with OTP
//	@Tags         accounts user
//	@Accept       json
//	@Produce      json
//	@Param        payload body model.RegisterInput true "payload"
//	@Success      200  {object}   response.ErrorResponseData
//	@Failure      400  {object}  response.ErrorResponseData
//	@Failure      500  {object}  response.ErrorResponseData
//	@Router      /user/register [post]
func (c *cUserLogin) Register(ctx *gin.Context) {
	var params model.RegisterInput
	if err := ctx.ShouldBind(&params); err != nil {
		fmt.Println("Error binding params: ", err)
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	codeStatus, err := service.UserLogin().Register(ctx, &params)
	if err != nil {
		global.Logger.Error("Error register user OTP: ", zap.Error(err))
		response.ErrorResponse(ctx, codeStatus, err.Error())
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}
