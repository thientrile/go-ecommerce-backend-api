package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func CheckValidParams(ctx *gin.Context, params interface{}) bool {
	if err := ctx.ShouldBind(&params); err != nil {
		fmt.Println("Error binding params: ", err)
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return false
	}
	return true
}

func GenerateNickname() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("user_%06d", rand.Intn(1000000))
}

// generate uuid
func GenerateUUID(userId int) string {
	newUUID := uuid.New()
	text := fmt.Sprintf("*%d*", userId)
	// Convert UUID to string
	uuidString := strings.ReplaceAll(newUUID.String(), "-", text)
	return strconv.Itoa(userId) + "_" + uuidString
}
