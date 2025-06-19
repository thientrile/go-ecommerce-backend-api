package impl

import (
	"context"
	"database/sql"
	"encoding/json"
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
	"go-ecommerce-backend-api.com/pkg/utils/auth"
	"go-ecommerce-backend-api.com/pkg/utils/crypto"
	"go-ecommerce-backend-api.com/pkg/utils/random"
	"go-ecommerce-backend-api.com/pkg/utils/sendto"
	"go.uber.org/zap"
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

// ---- Two-Factor Authentication ----
func (s *sUserLogin) IsTwoFactorEnabled(ctx context.Context, userId int64) (codeStatus int, rs bool, err error) {
	// Check if two-factor authentication is enabled for the user

	return response.ErrCodeSuccess, true, nil
}

func (s *sUserLogin) SetupTwoFactorAuth(ctx context.Context, in *model.SetupTwoFactorAuthInput) (codeStatus int, err error) {
	// logic
	// 1 check is two-factor authentication already enabled for the user
	isTwoFacorAuth, err := s.r.IsTwoFactorEnabled(ctx, int32(in.UserId))
	if err != nil {
		return response.ErrCodeTwoFactorAuthFailed, err
	}
	if isTwoFacorAuth > 0 {
		return response.ErrCodeTwoFactorAuthFailed, fmt.Errorf("two-factor authentication is already enabled for this user")
	}
	// 2 new type authe
	var twoFactorAuthType database.PreGoAccUserTwoFactor9999TwoFactorAuthType
	switch in.TwoFactorAuthType {
	case consts.EMAIL:
		twoFactorAuthType = database.PreGoAccUserTwoFactor9999TwoFactorAuthTypeEMAIL
		err = s.r.EnableTwoFactorTypeEmail(ctx, database.EnableTwoFactorTypeEmailParams{
			UserID:            int32(in.UserId),
			TwoFactorAuthType: twoFactorAuthType,
			TwoFactorEmail:    sql.NullString{String: in.TwoFactorEmail, Valid: true},
		})
	case consts.MOBILE:
		twoFactorAuthType = database.PreGoAccUserTwoFactor9999TwoFactorAuthTypeSMS
		err = s.r.EnableTwoFactorTypeSMS(ctx, database.EnableTwoFactorTypeSMSParams{
			UserID:            int32(in.UserId),
			TwoFactorAuthType: twoFactorAuthType,
			TwoFactorPhone:    sql.NullString{String: in.TwoFactorEmail, Valid: true},
		})
	default:
		twoFactorAuthType = database.PreGoAccUserTwoFactor9999TwoFactorAuthTypeAPP
		err = s.r.EnableTwoFactorTypeEmail(ctx, database.EnableTwoFactorTypeEmailParams{
			UserID:            int32(in.UserId),
			TwoFactorAuthType: twoFactorAuthType,
			TwoFactorEmail:    sql.NullString{String: in.TwoFactorEmail, Valid: true},
		})
	}

	if err != nil {
		global.Logger.Error("Error enabling two-factor authentication: ", zap.Error(err))
		return response.ErrCodeTwoFactorAuthFailed, err
	}
	// 3 send otp to in.TwoFactorEmail
	KeyUserTwoFator := crypto.GetHash("2fa:" + strconv.Itoa(int(in.UserId)))
	otpNew := random.GenerateSixDigitOTP()
	fmt.Printf("Generated OTP for user %d: %d\n", in.UserId, otpNew)
	go global.RDB.Set(ctx, KeyUserTwoFator, otpNew, time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
	switch in.TwoFactorAuthType {
	case consts.EMAIL:
		{

			err = sendto.SendTemplateEmailOtp([]string{in.TwoFactorEmail}, "thientrile2003@gmail.com", "otp-auth.html", map[string]interface{}{
				"otp":                strconv.Itoa(otpNew),
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

	default:
		{
			body := make(map[string]interface{})
			body["type"] = in.TwoFactorAuthType
			body["otp"] = strconv.Itoa(otpNew)
			body["email"] = in.TwoFactorEmail
			bodyRequest, _ := json.Marshal(body)
			err = sendto.SendMessageToKafka("Enable-auth-fator", string(bodyRequest))
			if err != nil {
				fmt.Printf("Error sending message to Kafka: %v\n", err)
				return response.ErrCodeSendEmailOtp, err
			}
		}
	}
	return response.ErrCodeSuccess, nil
}

func (s *sUserLogin) VerifyTwoFactorAuth(ctx context.Context, in *model.TwoFactorVerifycationInput) (codeStatus int, err error) {
	// Implement logic to verify the two-factor authentication code provided by the user
	// This could involve checking the code against the secret key stored in the database
	// For example, you might use a library like "github.com/pquerna/otp" to verify a TOTP code
	return response.ErrCodeSuccess, nil
}

// ---- End of Two-Factor Authentication ----
// Implement the methods of IUserLogin interface here

func (s *sUserLogin) Login(ctx context.Context, in *model.LoginInput) (codeStatus int, out model.LoginOutput, err error) {
	// Implement login logic
	// check if user exists in user base
	userBaseFound, err := s.r.GetOneUserInfo(ctx, in.Username)
	if err != nil {
		return response.ErrCodeAuthenticationFailed, out, nil
	}
	// check password
	if !crypto.MatchingPassword(userBaseFound.UserPassword, in.Password, userBaseFound.UserSalt) {
		return response.ErrCodeAuthenticationFailed, out, fmt.Errorf("does not match password")
	}
	// 3 check two factor authentication

	//4. update password time
	go s.r.LoginUserBase(ctx, database.LoginUserBaseParams{
		UserAccount:  in.Username,
		UserPassword: in.Password,
		UserLoginIp:  sql.NullString{String: "127.0.0.1", Valid: true},
	})
	// 5 create UUID
	subToken := utils.GenerateUUID(int(userBaseFound.UserID))
	// 6. get user info
	userInfo, err := s.r.GetUser(ctx, uint64(userBaseFound.UserID))
	if err != nil {
		global.Logger.Error("Error getting user info: ", zap.Error(err))
		return response.ErrCodeAuthenticationFailed, out, fmt.Errorf("failed to get user info: %v", err)
	}
	// 7. conver to json
	infoUserJson, err := json.Marshal(userInfo)
	if err != nil {

		return response.ErrCodeAuthenticationFailed, out, fmt.Errorf("failed to marshal user info: %v", err)
	}
	// 8. give infoUser to redis with key =subToken
	err = global.RDB.Set(ctx, subToken, infoUserJson, time.Duration(consts.TIME_LOGIN_LIFE)*time.Minute).Err()
	if err != nil {
		return response.ErrCodeAuthenticationFailed, out, fmt.Errorf("failed to set user info in Redis: %v", err)
	}
	// 9. create token
	tokenObj, err := auth.CreateToken(subToken)
	if err != nil {
		return response.ErrCodeAuthenticationFailed, out, err
	}
	out.Token.AccessToken = tokenObj.AccessToken
	out.Token.RefreshToken = tokenObj.RefreshToken
	out.Message = "Login successful"
	return response.ErrCodeSuccess, out, nil
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
		return response.ErrCodeUserHasExist, fmt.Errorf("user has already registered")
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

	// 4. check if otp exists in database
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

	default:
		{
			body := make(map[string]interface{})
			body["type"] = in.VerifyType
			body["otp"] = strconv.Itoa(optNew)
			body["email"] = in.VerifyKey
			bodyRequest, _ := json.Marshal(body)
			err = sendto.SendMessageToKafka("otp-auth", string(bodyRequest))
			if err != nil {
				fmt.Printf("Error sending message to Kafka: %v\n", err)
				return response.ErrCodeSendEmailOtp, err
			}
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
func (s *sUserLogin) VerifyOTP(ctx context.Context, in *model.VerifyInput) (out model.VerifyOtpOutput, err error) {
	// Implement login logic

	//check if user exists in redis
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))
	userKey := utils.GetUserKey(hashKey)
	// get otp
	otpFound, err := global.RDB.Get(ctx, userKey).Result()
	if err != nil {
		return out, fmt.Errorf("OTP not found or expired, please register again")
	}
	countKey := userKey + "_count"
	if in.VerifyCode != otpFound {
		// nếu otp không đúng thì tăng số lần thử MAX_COUNT_OTP
		count, err := global.RDB.Incr(ctx, countKey).Result()
		if err != nil {
			return out, fmt.Errorf("failed to increment OTP attempt count: %v", err)
		}
		if count > consts.MAX_COUNT_OTP {
			// If the count exceeds the maximum allowed attempts, delete the OTP key
			global.RDB.Del(ctx, userKey)
			return out, fmt.Errorf("maximum attempts exceeded, please register again")
		}
		if count == 1 {
			err = global.RDB.Expire(ctx, countKey, time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
			if err != nil {
				return out, fmt.Errorf("failed to set expiration for OTP: %v", err)
			}
		}
		if err != nil {
			return out, fmt.Errorf("failed to set expiration for OTP attempt count: %v", err)
		}
		return out, fmt.Errorf("invalid OTP code, please try again")

	}
	global.RDB.Del(ctx, userKey)
	// global.RDB.Del(ctx, countKey) // Reset the count if OTP is correct
	fmt.Printf("OTP found: %s\n", hashKey)
	infoOTP, err := s.r.GetValidOTP(ctx, hashKey)
	if err != nil {
		return out, fmt.Errorf("OTP invalid: %v", err)
	}
	// update status verify otp
	err = s.r.UpdateUserVerificationStatus(ctx, hashKey)
	if err != nil {
		return out, fmt.Errorf("OTP verification failed: %v", err)
	}
	//output
	out.Token = infoOTP.VerifyKeyHash
	out.Message = "OTP verified successfully"

	return out, err
}

func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context, in *model.UpdatePasswordRegisterInput) (codeStatus int, err error) {
	// Implement login logic
	infoOTP, err := s.r.GetInfoOTP(ctx, in.Token)
	fmt.Printf("infoOTP: %+v\n", infoOTP)
	if err != nil {
		return response.ErrCodeUserOtpNotExist, fmt.Errorf("OTP does not exist or is invalid: %v", err)
	}
	if infoOTP.IsVerified.Int32 == 0 {
		return response.ErrCodeUserOtpNotExist, fmt.Errorf("OTP exists but is not valid, please register again")
	}
	// update user base
	userBase := database.AddUserBaseParams{}
	userBase.UserAccount = infoOTP.VerifyKey
	userBase.UserPassword = crypto.GetHash(in.Password)
	userSalt, err := crypto.GenerateSalt(16)
	if err != nil {
		return response.ErrCodeUserOtpNotExist, err
	}
	userBase.UserSalt = userSalt
	userBase.UserPassword = crypto.HashPassword(in.Password, userSalt)
	// add userBase to user_base table
	NewUserBase, err := s.r.AddUserBase(ctx, userBase)
	if err != nil {
		return response.ErrCodeUserHasExist, fmt.Errorf("failed to add user base: %v", err)
	}
	lastIdVerify, err := NewUserBase.LastInsertId()
	if err != nil {
		return response.ErrCodeUserHasExist, err
	}
	userInfo := database.AddUserHaveUserIdParams{
		UserID:               uint64(lastIdVerify),
		UserAccount:          infoOTP.VerifyKey,
		UserNickname:         sql.NullString{String: utils.GenerateNickname(), Valid: true},
		UserAvatar:           sql.NullString{String: consts.USER_AVATAR_DEFAULT, Valid: true},
		UserState:            1,
		UserGender:           sql.NullInt16{Int16: 0, Valid: true},
		UserBirthday:         sql.NullTime{Time: time.Now(), Valid: true},
		UserIsAuthentication: uint8(1), // 1 for authenticated
	}
	switch infoOTP.VerifyType.Int32 {
	case int32(consts.EMAIL):
		{

			userInfo.UserEmail = sql.NullString{String: infoOTP.VerifyKey, Valid: true}
		}
	case int32(consts.MOBILE):
		{
			userInfo.UserMobile = sql.NullString{String: infoOTP.VerifyKey, Valid: true}
		}
	}
	// update user info by id
	_, err = s.r.AddUserHaveUserId(ctx, userInfo)
	if err != nil {

		return response.ErrCodeUserHasExist, err
	}
	return response.ErrCodeSuccess, nil
}
