package model

type RegisterInput struct {
	VerifyKey     string `json:"verify_key" binding:"required"`
	VerifyType    int    `json:"verify_type" binding:"required"`
	VerifyPurpose string `json:"verify_purpose" binding:"required"`
}
