package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-ecommerce-backend-api.com/internal/repo"
	"go-ecommerce-backend-api.com/pkg/response"
	"go-ecommerce-backend-api.com/pkg/utils/crypto"
	"go-ecommerce-backend-api.com/pkg/utils/random"
	"go-ecommerce-backend-api.com/pkg/utils/sendto"
)

// import "go-ecommerce-backend-api.com/internal/repo"

// type UserService struct {
// 	userRepo *repo.UserRepo
// }

// func NewUserService() *UserService {
// 	return &UserService{
// 		userRepo: repo.NewUserRepo(),
// 	}
// }
// func (us *UserService) GetInfoUser() string {
// 	return us.userRepo.GetInfoUser()
// }

// interface UserServiceInterface

type IUserService interface {
	Regisger(email string, purpose string) int
}

type userService struct {
	userRepo     repo.IUserRepository
	userAuthRepo repo.IUserAuthRepository
	//..
}

func NewUserService(
	userRepo repo.IUserRepository,
	userAuthRepo repo.IUserAuthRepository,
) IUserService {
	return &userService{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}

// Regisger implements IUserService.
func (us *userService) Regisger(email string, purpose string) int {
	// 0.hash email
	hashEmail := crypto.GetHash(email)
	fmt.Printf("Hash Email: %s\n", hashEmail)
	// 5. check OTP is available in Redis

	// 6.user spams

	//  1. check email exist in database
	if us.userRepo.GetUserByEmail(email) {
		return response.ErrCodeUserHasExist
	}
	// 2. new OTP ...
	otp := random.GenerateSixDigitOTP()
	if purpose == "TEXT_USER" {
		otp = 123456
	}
	fmt.Printf("OTP::: %d\n", otp)

	// 3. save OTP in Redis with expiration time
	expirationMinutes := 10 * time.Minute
	err := us.userAuthRepo.AddOTP(hashEmail, otp, int64(expirationMinutes)) // 5 minutes expiration
	if err != nil {
		return response.ErrInvalidOTP
	}
	// 4. send Email OTP
	// err = sendto.SendTemplateEmailOtp([]string{email}, "thientrile2003@gmail.com", "otp-auth.html", map[string]interface{}{
	// 	"otp":                strconv.Itoa(otp),
	// 	"expiration_minutes": expirationMinutes,
	// })
	// if err != nil {
	// 	return response.ErrCodeSendEmailOtp
	// }
	// send OTP via Kafka
	body := make (map[string]interface{})
	body["otp"] = strconv.Itoa(otp)
	body["email"] = email
	bodyRequest, _ := json.Marshal(body)
	err = sendto.SendMessageToKafka("otp-auth", string(bodyRequest))
	if err != nil {
		fmt.Printf("Error sending message to Kafka: %v\n", err)
		return response.ErrCodeSendEmailOtp
	}
	return response.ErrCodeSuccess
}
