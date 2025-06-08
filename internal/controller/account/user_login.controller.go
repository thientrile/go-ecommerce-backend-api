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

// manager controller login user
var Login = new(cUserLogin)

// cUserLogin is the controller for user login related actions.
type cUserLogin struct{}

// UpdatePasswordRegister documentation
//
// @Summary      Update password after registration
// @Description  Allows a user to set or update their password after registration
// @Tags         accounts user
// @Accept       json
// @Produce      json
// @Param        payload body model.UpdatePasswordRegisterInput true "payload"
// @Success      200  {object}   response.ErrorResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/update-password-register [post]
func (c *cUserLogin) UpdatePasswordRegister(ctx *gin.Context) {
	var params model.UpdatePasswordRegisterInput
	if !utils.CheckValidParams(ctx, &params) {
		return
	}
	codeStatus, err := service.UserLogin().UpdatePasswordRegister(ctx, &params)
	if err != nil {
		global.Logger.Error("Error updating password after registration: ", zap.Error(err))
		response.ErrorResponse(ctx, codeStatus, "")
		return
	}
	response.SuccessResponse(ctx, codeStatus, nil)
}

/*
User Login documentation

@Summary      User login with credentials
@Description  Login user and return authentication token
@Tags         accounts user
@Accept       json
@Produce      json
@Param        payload body model.LoginInput true "payload"
@Success      200  {object}   response.ErrorResponseData
@Failure      400  {object}  response.ErrorResponseData
@Failure      500  {object}  response.ErrorResponseData
@Router       /user/login [post]
*/
func (c *cUserLogin) Login(ctx *gin.Context) {
	// Implement login logic
	err := service.UserLogin().Login(ctx)
	if err != nil {

		response.ErrorResponse(ctx, response.ErrInvalidOTP, err.Error())
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)

}

// Verify OTP documentation
//
// @Summary      Verify OTP for user
// @Description  Verify the OTP sent to the user
// @Tags         accounts user
// @Accept       json
// @Produce      json
// @Param        payload body model.VerifyInput true "payload"
// @Success      200  {object}   response.ErrorResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/verify-otp [post]
func (c *cUserLogin) VerifyOTP(ctx *gin.Context) {
	var params model.VerifyInput
	if !utils.CheckValidParams(ctx, &params) {
		return
	}

	payload, err := service.UserLogin().VerifyOTP(ctx, &params)
	if err != nil {
		global.Logger.Error("Error verifying OTP: ", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrInvalidOTP, "")
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, payload)
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
	if !utils.CheckValidParams(ctx, &params) {
		return
	}
	codeStatus, err := service.UserLogin().Register(ctx, &params)
	if err != nil {
		global.Logger.Error("Error register user OTP: ", zap.Error(err))
		response.ErrorResponse(ctx, codeStatus, "")
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}
