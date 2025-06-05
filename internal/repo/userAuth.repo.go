package repo

import (
	"fmt"
	"time"

	"go-ecommerce-backend-api.com/global"
)

type IUserAuthRepository interface {
	AddOTP(email string, otp int, expirationTime int64) error
}

// userAuthRepository is a concrete implementation of IUserAuthRepository
type userAuthRepository struct{}

// AddOTP implements the IUserAuthRepository interface.
func (r *userAuthRepository) AddOTP(email string, otp int, expirationTime int64) error {
	// TODO: implement the logic to add OTP
	key := fmt.Sprintf("usr:%s:otp", email) // user:email:otp
	return global.RDB.SetEx(ctx, key, otp, time.Duration(expirationTime)).Err()
}

func NewUserAuthRepository() IUserAuthRepository {
	return &userAuthRepository{}
}
