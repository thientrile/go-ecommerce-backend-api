package account

import (
	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/internal/model"
	"go-ecommerce-backend-api.com/internal/service"
	"go-ecommerce-backend-api.com/pkg/response"
	"go-ecommerce-backend-api.com/pkg/utils"
	"go.uber.org/zap"
)

type cUser2FA struct{}

var TwoFA = new(cUser2FA)

// Setup Two-Factor Authentication documentation
//
// @Summary      Setup two-factor authentication for user
// @Description  Enables two-factor authentication for the user account
// @Tags         accounts 2fa
// @Param 	  Authorization header string true "Bearer{token}"
// @Accept       json
// @Produce      json
// @Param        payload body model.SetupTwoFactorAuthInput true "payload"
// @Success      200  {object}   response.ErrorResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/setup-2fa [post]
func (c *cUser2FA) SetupTwoFactorAuth(ctx *gin.Context) {
	var params model.SetupTwoFactorAuthInput
	if !utils.CheckValidParams(ctx, &params) {
		return
	}
	codeStatus, err := service.UserLogin().SetupTwoFactorAuth(ctx, &params)
	if err != nil {
		global.Logger.Error("Error setting up two-factor authentication: ", zap.Error(err))
		response.ErrorResponse(ctx, codeStatus, "")
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}
//
// Verify Two-Factor Authentication documentation
//
// @Summary      Verify two-factor authentication code
// @Description  Verifies the two-factor authentication code for the user account
// @Tags         accounts 2fa
// @Param        Authorization header string true "Bearer{token}"
// @Accept       json
// @Produce      json
// @Param        payload body model.TwoFactorVerifycationInput true "payload"
// @Success      200  {object}   response.ErrorResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/verify-2fa [post]
func (c *cUser2FA) VerifyTwoFactorAuth(ctx *gin.Context) {

	var params model.TwoFactorVerifycationInput
	if !utils.CheckValidParams(ctx, &params) {
		return
	}
	codeStatus, err := service.UserLogin().VerifyTwoFactorAuth(ctx, &params)
	if err != nil {
		global.Logger.Error("Error verifying two-factor authentication: ", zap.Error(err))
		response.ErrorResponse(ctx, codeStatus, "")
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}
