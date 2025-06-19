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

	"github.com/redis/go-redis/v9"
	"go-ecommerce-backend-api.com/global"
	consts "go-ecommerce-backend-api.com/internal/const"
	"go-ecommerce-backend-api.com/internal/database"
	"go-ecommerce-backend-api.com/internal/model"
	pkgContext "go-ecommerce-backend-api.com/pkg/context"
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
	// 0 check userId is valid
	userId, err := pkgContext.GetUserIdFormUUID(ctx)
	if err != nil {
		global.Logger.Error("Error getting user ID from context: ", zap.Error(err))
		return response.ErrCodeTwoFactorAuthFailed, err
	}
	hashKey := crypto.GetHash(strconv.FormatUint(userId, 10))
	KeyUserTwoFator := utils.GetUserKey(hashKey)
	// 1 check is two-factor authentication already enabled for the user
	isTwoFacorAuth, err := s.r.IsTwoFactorEnabled(ctx, int32(userId))
	if err != nil {
		return response.ErrCodeTwoFactorAuthFailed, err
	}
	if isTwoFacorAuth > 0 {
		return response.ErrCodeTwoFactorAuthFailed, fmt.Errorf("two-factor authentication is already enabled for this user")
	}
	// 2 new type authentication
	twoFactorAuthType := database.PreGoAccUserTwoFactor9999TwoFactorAuthType(getAuthType(in.TwoFactorAuthType))
	switch in.TwoFactorAuthType {
	case consts.EMAIL:
		{
			err = s.r.EnableTwoFactorTypeEmail(ctx, database.EnableTwoFactorTypeEmailParams{
				UserID:            int32(userId),
				TwoFactorAuthType: twoFactorAuthType,
				TwoFactorEmail:    sql.NullString{String: in.TwoFactorEmail, Valid: true},
			})
		}
	default:
		{
			err = s.r.EnableTwoFactorTypeSMS(ctx, database.EnableTwoFactorTypeSMSParams{
				UserID:            int32(userId),
				TwoFactorAuthType: twoFactorAuthType,
				TwoFactorPhone:    sql.NullString{String: in.TwoFactorEmail, Valid: true}, // Assuming TwoFactorEmail is used for phone in this case
			})
		}
	}
	if err != nil {
		global.Logger.Error("Error enabling two-factor authentication: ", zap.Error(err))
		return response.ErrCodeTwoFactorAuthFailed, err
	}
	// 3 send otp to in.TwoFactorEmail

	otpNew := random.GenerateSixDigitOTP()
	fmt.Printf("Generated OTP for user %d: %d\n", userId, otpNew)
	go global.RDB.SetEx(ctx, KeyUserTwoFator, otpNew, time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
	err = sendOtp(in.TwoFactorEmail, in.TwoFactorAuthType, strconv.Itoa(otpNew))
	if err != nil {
		global.Logger.Error("Error sending OTP: ", zap.Error(err))
		return response.ErrCodeSendEmailOtp, fmt.Errorf("failed to send OTP: %v", err)
	}
	return response.ErrCodeSuccess, nil
}

func (s *sUserLogin) VerifyTwoFactorAuth(ctx context.Context, in *model.TwoFactorVerifycationInput) (codeStatus int, err error) {
	// 0 check userId is valid
	userId, err := pkgContext.GetUserIdFormUUID(ctx)
	if err != nil {
		global.Logger.Error("Error getting user ID from context: ", zap.Error(err))
		return response.ErrCodeTwoFactorAuthFailed, err
	}
	hashKey := crypto.GetHash(strconv.FormatUint(userId, 10))
	KeyUserTwoFator := utils.GetUserKey(hashKey)

	//1. check is two-factor authentication enabled for the user
	isTwoFacorAuth, err := s.r.IsTwoFactorEnabled(ctx, int32(userId))
	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailded, err
	}
	if isTwoFacorAuth > 0 {
		return response.ErrCodeTwoFactorAuthVerifyFailded, fmt.Errorf("two-factor authentication is not enabled for this user")
	}

	// 2. validate OTP
	if err := validateOTP(ctx, KeyUserTwoFator, in.TwoFactorCode); err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailded, err
	}

	//4 update status two-factor authentication
	TwoFactorAuthType := database.PreGoAccUserTwoFactor9999TwoFactorAuthType(getAuthType(in.TwoFactorAuthType))
	err = s.r.UpdateTwoFactorStatusVerification(ctx, database.UpdateTwoFactorStatusVerificationParams{
		UserID:            int32(userId),
		TwoFactorAuthType: TwoFactorAuthType,
	})
	if err != nil {
		global.Logger.Error("Error updating two-factor authentication status: ", zap.Error(err))
		return response.ErrCodeTwoFactorAuthVerifyFailded, fmt.Errorf("failed to update two-factor authentication status: %v", err)
	}
	return response.ErrCodeSuccess, nil
}

// validateOTP handles OTP validation and attempt counting for two-factor authentication.

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
	isTwoFactorEnabled, err := s.r.IsTwoFactorEnabled(ctx, int32(userBaseFound.UserID))
	if err != nil {
		global.Logger.Error("Error checking two-factor authentication status: ", zap.Error(err))
		return response.ErrCodeAuthenticationFailed, out, fmt.Errorf("failed to check two-factor authentication status: %v", err)
	}
	if isTwoFactorEnabled > 0 {
		// found two-factor authentication is enabled for this user
		foundFactorMethods, err := s.r.GetUserFactorMethods(ctx, int32(userBaseFound.UserID))
		if err != nil {
			global.Logger.Error("Error getting two-factor method: ", zap.Error(err))
			return response.ErrCodeAuthenticationFailed, out, fmt.Errorf("failed to get two-factor method: %v", err)
		}
		if len(foundFactorMethods) == 0 {
			global.Logger.Error("No two-factor methods found for user")
			return response.ErrCodeTwoFactorEnabled, out, fmt.Errorf("no two-factor methods found for user")
		}
		authType := convertAuthType(string(foundFactorMethods[0].TwoFactorAuthType))
		// set OTP for two-factor authentication
		hashKey := crypto.GetHash(strconv.FormatUint(uint64(userBaseFound.UserID), 10))
		KeyUserTwoFator := utils.GetUserKey(hashKey)
		otpNew := random.GenerateSixDigitOTP()
		go global.RDB.SetEx(ctx, KeyUserTwoFator, otpNew, time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
		// send OTP to user
		key := map[int]string{
			consts.EMAIL: foundFactorMethods[0].TwoFactorEmail.String,
			consts.SMS:   foundFactorMethods[0].TwoFactorPhone.String,
		}
		go sendOtp(key[authType], authType, strconv.Itoa(otpNew))
		out.Message = fmt.Sprintf("Two-factor authentication is enabled. Please verify your OTP sent to %s", foundFactorMethods[0].TwoFactorAuthType)
		return response.ErrCodeSuccess, out, nil
	}
	//get ip address from context
	ipClient := ctx.Value("IpAddress")
	ipStr, _ := ipClient.(string)
	//4. update password time

	go s.r.LoginUserBase(ctx, database.LoginUserBaseParams{
		UserAccount:  in.Username,
		UserPassword: in.Password,
		UserLoginIp:  sql.NullString{String: ipStr, Valid: true},
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

// Register handles user registration and OTP verification.
// It generates an OTP, saves it in Redis, and sends it to the user's email.
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
	if err = sendOtp(in.VerifyKey, in.VerifyType, strconv.Itoa(optNew)); err != nil {
		return response.ErrCodeSendEmailOtp, err
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
	if err := validateOTP(ctx, userKey, in.VerifyCode); err != nil {
		return out, fmt.Errorf("OTP verification failed: %v", err)
	}

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
	case int32(consts.SMS):
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

// validateOTP validates the OTP for two-factor authentication and manages attempt counting.

func validateOTP(ctx context.Context, key, inputCode string) error {
	otpVerifyAuth, err := global.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		global.Logger.Error("Error getting OTP from Redis: ", zap.Error(err))
		return fmt.Errorf("OTP not found or expired, please try again")
	} else if err != nil {
		return fmt.Errorf("OTP not found or expired, please try again")
	}
	countKey := key + "_count"
	if inputCode != otpVerifyAuth {
		count, err := global.RDB.Incr(ctx, countKey).Result()
		if err != nil {
			return fmt.Errorf("failed to increment OTP attempt count: %v", err)
		}
		if count > consts.MAX_COUNT_OTP {
			global.RDB.Del(ctx, key)
			return fmt.Errorf("maximum attempts exceeded, please verify again")
		}
		if count == 1 {
			err = global.RDB.Expire(ctx, countKey, time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
			if err != nil {
				return fmt.Errorf("failed to set expiration for OTP: %v", err)
			}
		}
		if err != nil {
			return fmt.Errorf("failed to set expiration for OTP attempt count: %v", err)
		}
		return fmt.Errorf("invalid OTP code, please try again")
	}
	// If OTP is valid, delete the OTP key and reset the attempt count
	global.RDB.Del(ctx, key)
	return nil
}

// getAuthType converts the integer authType to a string representation.
// It returns "EMAIL", "SMS", "APP", or "UNKNOWN" if the auth
func getAuthType(authType int) string {
	switch authType {
	case 1:
		return "EMAIL"
	case 2:
		return "SMS"
	case 3:
		return "APP"
	default:
		return "UNKNOWN"
	}
}
func convertAuthType(authType string) int {
	switch strings.ToUpper(authType) {
	case "EMAIL":
		return consts.EMAIL
	case "SMS":
		return consts.SMS
	case "APP":
		return consts.APP
	default:
		return 0
	}
}

// sendOtp sends an OTP to the specified authKey based on the authType.
// It supports sending OTPs via email, SMS, or Kafka message.
func sendOtp(authKey string, authType int, otpNew string) (err error) {
	switch authType {
	case consts.EMAIL:
		{

			err = sendto.SendTemplateEmailOtp([]string{authKey}, "thientrile2003@gmail.com", "otp-auth.html", map[string]interface{}{
				"otp":                otpNew,
				"expiration_minutes": consts.TIME_OTP_REGISTER,
			})
			if err != nil {
				return err
			}

		}
	case consts.SMS:
		{
			fmt.Printf("Sending OTP %s to phone number %s\n", otpNew, authKey)
			return nil
		}

	default:
		{
			body := make(map[string]interface{})
			body["type"] = authType
			body["otp"] = otpNew
			body["key"] = authKey
			bodyRequest, _ := json.Marshal(body)
			err = sendto.SendMessageToKafka("otp-auth", string(bodyRequest))
			if err != nil {
				fmt.Printf("Error sending message to Kafka: %v\n", err)
				return err
			}
		}
	}
	return nil
}
