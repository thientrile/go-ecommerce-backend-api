package response

const (
	ErrCodeSuccess      = 20001 // Success
	ErrCodeParamInvalid = 20003 // Parameter Invalid
	ErrCodeTokenInvalid = 30001 // Token Invalid
	ErrInvalidOTP       = 30002 // Invalid OTP
	ErrCodeSendEmailOtp = 30003 // Failed to send email OTP
	// register
	ErrCodeUserHasExist = 50001 // User has exist
	// login
	ErrCodeOtpNotExist      = 60000 // OTP does not exist
	ErrCodeUserOtpNotExist  = 60008 // OTP exists but is not valid
	ErrCodeUserNotExist     = 60009 // User does not exist
	ErrCodeTwoFactorEnabled = 60012 // Two-factor authentication is enabled
	ErrCodeOTPExisted       = 60013 // OTP already exists for user
	//authentication
	ErrCodeAuthenticationFailed  = 60010 // Authentication failed
	ErrCodeAuthenticationSuccess = 60011 // Authentication success
	// two factor authentication
	ErrCodeTwoFactorAuthFailed        = 80001 // Two factor authentication failed
	ErrCodeTwoFactorAuthVerifyFailded = 80002 // Two factor authentication success
	// rate limit
	ErrCodeRateLimitExceeded = 90001 // Rate limit exceeded
	ErrCodeRateLimitNotFound = 90002 // Rate limit not found
	ErrCodeRateLimitError    = 90003 // Rate limit error
	// Ticket
	ErrCodeTicketItemNotFound = 100001 // Ticket item not found
	ErrCodeTicketItemError    = 100002 // Ticket item error
)

var msg = map[int]string{
	ErrCodeSuccess:      "Success",
	ErrCodeParamInvalid: "Parameter Invalid",
	ErrCodeTokenInvalid: "Token Invalid",
	ErrInvalidOTP:       "Invalid OTP",
	ErrCodeUserHasExist: "User has exist",
	ErrCodeSendEmailOtp: "Failed to send email OTP",
	// login
	ErrCodeOtpNotExist:      "OTP exists but is not valid",
	ErrCodeUserOtpNotExist:  "OTP does not exist",
	ErrCodeUserNotExist:     "User does not exist",
	ErrCodeTwoFactorEnabled: "Two-factor authentication is enabled",
	// authentication
	ErrCodeAuthenticationFailed:  "Authentication failed",
	ErrCodeAuthenticationSuccess: "Authentication success",
	// two factor authentication
	ErrCodeTwoFactorAuthVerifyFailded: "Two factor authentication verification failed",

	// rate limit
	ErrCodeRateLimitExceeded: "Rate limit exceeded",
	ErrCodeRateLimitNotFound: "Rate limit not found",
	ErrCodeRateLimitError:    "Rate limit error",
	ErrCodeOTPExisted:        "OTP already exists for user",
	// Ticket
	ErrCodeTicketItemNotFound: "Ticket item not found",
	ErrCodeTicketItemError:    "Ticket item error",
}
