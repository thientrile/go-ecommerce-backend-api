package account

import (
	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/internal/service"
	"go-ecommerce-backend-api.com/pkg/response"
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
