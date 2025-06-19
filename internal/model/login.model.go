package model

type RegisterInput struct {
	VerifyKey     string `json:"verify_key" binding:"required"`
	VerifyType    int    `json:"verify_type" binding:"required"`
	VerifyPurpose string `json:"verify_purpose" binding:"required"`
}
type VerifyInput struct {
	VerifyKey  string `json:"verify_key" binding:"required"`
	VerifyCode string `json:"verify_code" binding:"required"`
}
type VerifyOtpOutput struct {
	Token string `json:"token"`
	// UserId  int64  `json:"user_id"`
	Message string `json:"message"`
}
type UpdatePasswordRegisterInput struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginOutput struct {
	Token struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	Message string `json:"message"`
}

// two factor authentication
type SetupTwoFactorAuthInput struct {
	UserId            uint32 `json:"user_id" binding:"required"`
	TwoFactorAuthType int    `json:"two_factor_auth_type" binding:"required"` // 1:"email" and 2:"sms"
	TwoFactorEmail    string `json:"two_factor_email,omitempty"`              // required if TwoFactorAuthType is "email"
}

type TwoFactorVerifycationInput struct {
	UserId            uint32 `json:"user_id" binding:"required"`
	TwoFactorCode     string `json:"two_factor_code" binding:"required"`      // the code to verify
	TwoFactorAuthType string `json:"two_factor_auth_type" binding:"required"` // "email" or "sms"
}
