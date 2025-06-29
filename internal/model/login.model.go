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
	Token   string `json:"token"`
	Message string `json:"message"`
}

// two factor authentication
type SetupTwoFactorAuthInput struct {
	TwoFactorAuthType int    `json:"two_factor_auth_type" binding:"required"` // 1:"email" and 2:"sms" or 3:"app"
	TwoFactorEmail    string `json:"two_factor_email,omitempty"`              // required if TwoFactorAuthType is "email"
}

type TwoFactorVerifycationInput struct {
	TwoFactorCode     string `json:"two_factor_code" binding:"required"`      // the code to verify
	TwoFactorAuthType int    `json:"two_factor_auth_type" binding:"required"` // 1:"email" and 2:"sms" or 3:"app"
}
type TwoFactorVerifyOtp struct {
	VerifyKey          string `json:"verify_key" binding:"required"`      // the key to verify
	TwoFactorCode      string `json:"two_factor_code" binding:"required"` // the code to verify
	TwoFactorAuthToken string `json:"two_factor_auth_token" binding:"required"`
}
