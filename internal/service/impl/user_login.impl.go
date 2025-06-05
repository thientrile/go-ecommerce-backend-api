package impl

import (
	"context"

	"go-ecommerce-backend-api.com/internal/database"
)

type sUserLogin struct {
	// Implement the IUserLogin interface here
	r *database.Queries
}

func NewUserLoginImpl(r *database.Queries) * sUserLogin {

	return &sUserLogin{
		r: r,
	}
}

// Implement the methods of IUserLogin interface here

func (s * sUserLogin) Login(ctx context.Context) error {
	// Implement login logic
	return nil
}


func (s * sUserLogin) Register(ctx context.Context) error {
	// Implement login logic
	return nil
}
func (s * sUserLogin) VerifyOTP(ctx context.Context) error {
	// Implement login logic
	return nil
}

func (s * sUserLogin) UpdatePasswordRegister(ctx context.Context) error {
	// Implement login logic
	return nil
}