package model

type RegisterInput struct {
	VerifyKey     string `json:"verify_key" binding:"required"`
	VerifyType    int    `json:"verify_type" binding:"required"`
	VerifyPurpose string `json:"verify_purpose" binding:"required"`
}
type VerifyInput struct {
	VerifyKey     string `json:"verify_key" binding:"required"`
	VerifyCode    string `json:"verify_code" binding:"required"`
}
type VerifyOtpOutput struct {
	
	Token string `json:"token"`
	UserId int64  `json:"user_id"`
	Message string `json:"message"`
}