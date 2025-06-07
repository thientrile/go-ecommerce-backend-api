package utils

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"go-ecommerce-backend-api.com/pkg/response"
)

func GetUserKey(hashKey string) string {
	return fmt.Sprintf("u:%s:otp", hashKey)
}

func HandleRedisGetOTPError(err error, otpFound string) (responseCode int, errResult error, result bool) {
	switch {
	case err == redis.Nil:
		fmt.Println("key does not exist, creating new OTP")
		// Cho phép tiếp tục xử lý tạo OTP
		return 0, err, true

	case err != nil:
		fmt.Println("get failed, err: ", err)
		return response.ErrInvalidOTP, err, false

	case otpFound != "":
		return response.ErrCodeOtpNotExist, fmt.Errorf("OTP is not valid"), false
	}

	return 0, nil, true
}
