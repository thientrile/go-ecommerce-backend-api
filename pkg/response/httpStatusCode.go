package response

const (
	ErrCodeSuccess      = 20001 // Success
	ErrCodeParamInvalid = 20003 // Parameter Invalid
	ErrCodeTokenInvalid = 30001 // Token Invalid
	ErrInvalidOTP		= 30002 // Invalid OTP
	ErrCodeSendEmailOtp = 30003 // Failed to send email OTP
	// register
	ErrCodeUserHasExist = 50001 // User has exist

)

var msg = map[int]string{
	ErrCodeSuccess:      "Success",
	ErrCodeParamInvalid: "Parameter Invalid",
	ErrCodeTokenInvalid: "Token Invalid",
	ErrInvalidOTP:       "Invalid OTP",
	ErrCodeUserHasExist: "User has exist",
	ErrCodeSendEmailOtp: "Failed to send email OTP",
}
