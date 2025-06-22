package service

import (
	"context"

	"go-ecommerce-backend-api.com/internal/model"
)

type (
	//... interfaces
	IUserLogin interface {
		Login(ctx context.Context, in *model.LoginInput) (codeStatus int, out model.LoginOutput, err error)
		Register(ctx context.Context, in *model.RegisterInput) (codeResult int, err error)
		VerifyOTP(ctx context.Context, in *model.VerifyInput) (out model.VerifyOtpOutput, err error)
		UpdatePasswordRegister(ctx context.Context, in *model.UpdatePasswordRegisterInput) (userId int, err error)

		//  two-factor authentication
		IsTwoFactorEnabled(ctx context.Context, userId int64) (codeStatus int, rs bool, err error)
		// setup two-factor authentication
		SetupTwoFactorAuth(ctx context.Context, in *model.SetupTwoFactorAuthInput) (codeStatus int, err error)

		// verify Two-Factor Authentication
		VerifyTwoFactorAuth(ctx context.Context, in *model.TwoFactorVerifycationInput) (codeStatus int, err error)
		// verify Two-Factor Authentication otp
		VerifyTwoFactorAuthOTP(ctx context.Context, in *model.TwoFactorVerifyOtp) (codeStatus int, out model.LoginOutput, err error)
	}
	IUserInfo interface {
		GetInfoByUserId(ctx context.Context) error
	}
	IUserAdmin interface {
		RemoveUser(ctx context.Context) error
		FindOneUser(ctx context.Context) error
	}
)

var (
	localUserAdmin IUserAdmin
	localUserInfo  IUserInfo
	localUserLogin IUserLogin
)

func UserAdmin() IUserAdmin {
	if localUserAdmin == nil {
		panic("implement localUserAdmin not found interface IUserAdmin")
	}
	return localUserAdmin
}

func InitUserAdmin(i IUserAdmin) {
	localUserAdmin = i
}

func UserInfo() IUserInfo {
	if localUserInfo == nil {
		panic("implement localUserInfo not found interface IUserInfo")
	}
	return localUserInfo
}

func InitUserInfo(i IUserInfo) {
	localUserInfo = i
}

func UserLogin() IUserLogin {
	if localUserLogin == nil {
		panic("implement localUserLogin not found interface IUserLogin")
	}
	return localUserLogin
}

func InitUserLogin(i IUserLogin) {
	localUserLogin = i
}
