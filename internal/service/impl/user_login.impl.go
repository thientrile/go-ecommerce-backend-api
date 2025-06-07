package impl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go-ecommerce-backend-api.com/global"
	consts "go-ecommerce-backend-api.com/internal/const"
	"go-ecommerce-backend-api.com/internal/database"
	"go-ecommerce-backend-api.com/internal/model"
	"go-ecommerce-backend-api.com/pkg/response"
	"go-ecommerce-backend-api.com/pkg/utils"
	"go-ecommerce-backend-api.com/pkg/utils/crypto"
	"go-ecommerce-backend-api.com/pkg/utils/random"
	"go-ecommerce-backend-api.com/pkg/utils/sendto"
)

type sUserLogin struct {
	// Implement the IUserLogin interface here
	r *database.Queries
}

func NewUserLoginImpl(r *database.Queries) *sUserLogin {

	return &sUserLogin{
		r: r,
	}
}

// Implement the methods of IUserLogin interface here

func (s *sUserLogin) Login(ctx context.Context) error {
	// Implement login logic
	return nil
}

func (s *sUserLogin) Register(ctx context.Context, in *model.RegisterInput) (codeResult int, err error) {
	// Implement login registration
	// 1. hash email
	fmt.Printf("verifyKey: %s\n", in.VerifyKey)
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))
	fmt.Printf("Hash email: %s\n", hashKey)

	// 2. check if email exists in user base

	userFound, err := s.r.CheckUserBaseExists(ctx, in.VerifyKey)
	if err != nil {
		return response.ErrCodeUserHasExist, err

	}
	if userFound > 0 {
		return response.ErrCodeUserHasExist, fmt.Errorf("User has already registered")
	}

	// create otp
	userKey := utils.GetUserKey(hashKey)
	otpFound, err := global.RDB.Get(ctx, userKey).Result()
	redisCode, err, check := utils.HandleRedisGetOTPError(err, otpFound)
	if !check {
		return redisCode, err
	}
	optNew := random.GenerateSixDigitOTP()
	if in.VerifyPurpose == "TEST_USER" {
		optNew = 123456 // For testing purposes, use a fixed OTP
	}
	fmt.Printf("Generated OTP: %d\n", optNew)
	// 3. save OTP in Redis with expiration time

	err = global.RDB.SetEx(ctx, userKey, strconv.Itoa(optNew), time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
	if err != nil {
		return response.ErrInvalidOTP, err
	}
	//6 Send OTP to email
	switch in.VerifyType {
	case consts.EMAIL:
		{

			err = sendto.SendTemplateEmailOtp([]string{in.VerifyKey}, "thientrile2003@gmail.com", "otp-auth.html", map[string]interface{}{
				"otp":                strconv.Itoa(optNew),
				"expiration_minutes": consts.TIME_OTP_REGISTER,
			})
			if err != nil {
				return response.ErrCodeSendEmailOtp, err
			}

		}
	case consts.MOBILE:
		{
			return response.ErrCodeSuccess, nil
		}
	}
	// 7. save user to database
	result, err := s.r.InsertOTPVerify(ctx, database.InsertOTPVerifyParams{
		VerifyOtp:     strconv.Itoa(optNew),
		VerifyType:    sql.NullInt32{Int32: 1, Valid: true},
		VerifyKey:     in.VerifyKey,
		VerifyKeyHash: hashKey,
	})
	if err != nil {
		return response.ErrCodeSendEmailOtp, err
	}
	// 8. getlast id
	lastIdVerify, err := result.LastInsertId()
	if err != nil {
		return response.ErrCodeSendEmailOtp, err
	}
	log.Printf("Last Insert ID: %d\n", lastIdVerify)

	return response.ErrCodeSuccess, nil
}
func (s *sUserLogin) VerifyOTP(ctx context.Context) error {
	// Implement login logic
	return nil
}

func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context) error {
	// Implement login logic
	return nil
}
